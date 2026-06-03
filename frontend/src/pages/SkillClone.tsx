import { useState } from 'react';
import {
  Button,
  Card,
  Input,
  Form,
  Tag,
  Spin,
  message,
  Space,
  Typography,
} from 'antd';
import {
  ExperimentOutlined,
  CheckCircleOutlined,
  LoadingOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import {
  createSkill,
  processSkill,
  getSkill,
  getTask,
} from '../api/endpoints';
import type { Task } from '../api/types';
import type { PreviewState } from '../App';

const { TextArea } = Input;
const { Title, Paragraph } = Typography;

const STATUS_MAP: Record<string, { color: string; icon: React.ReactNode; text: string }> = {
  pending: { color: 'default', icon: <LoadingOutlined />, text: '等待中' },
  processing: { color: 'processing', icon: <LoadingOutlined />, text: '处理中' },
  done: { color: 'success', icon: <CheckCircleOutlined />, text: '已完成' },
  completed: { color: 'success', icon: <CheckCircleOutlined />, text: '已完成' },
  failed: { color: 'error', icon: <CloseCircleOutlined />, text: '失败' },
};

interface Props {
  setPreview: (p: PreviewState) => void;
}

export default function SkillClone({ setPreview }: Props) {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<string>('');

  const showInPreview = (title: string, content: React.ReactNode) => {
    setPreview({ visible: true, title, content });
  };

  const pollTask = async (taskId: string) => {
    let attempts = 0;
    const maxAttempts = 60;
    while (attempts < maxAttempts) {
      await new Promise((r) => setTimeout(r, 2000));
      try {
        const { data: task } = await getTask(taskId) as { data: Task };
        if (task.status === 'done' || task.status === 'completed') {
          const { data: skill } = await getSkill(task.ref_id);
          setStatus('done');
          // 在右侧预览区显示结果
          showInPreview('🧠 生成的思维框架', (
            <pre style={{
              background: '#f6f8fa',
              padding: 16,
              borderRadius: 8,
              overflow: 'auto',
              maxHeight: 'calc(100vh - 240px)',
              fontSize: 13,
              lineHeight: 1.6,
              whiteSpace: 'pre-wrap',
            }}>
              {skill.result}
            </pre>
          ));
          message.success('思维框架生成成功！');
          return;
        }
        if (task.status === 'failed') {
          setStatus('failed');
          message.error(`处理失败: ${task.error || '未知错误'}`);
          return;
        }
        setStatus(task.status);
        // 处理中也更新预览
        showInPreview('⏳ 处理中...', (
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin size="large" />
            <div style={{ marginTop: 16, color: '#888' }}>
              AI 正在分析你的语料，生成思维框架…
            </div>
          </div>
        ));
      } catch {
        // continue polling
      }
      attempts++;
    }
  };

  const onFinish = async (values: { name: string; corpus: string }) => {
    setLoading(true);
    setStatus('pending');
    try {
      const { data: skill } = await createSkill(values);
      setStatus('processing');

      showInPreview('⏳ 处理中...', (
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Spin size="large" />
          <div style={{ marginTop: 16, color: '#888' }}>正在提交任务…</div>
        </div>
      ));

      const { data: task } = await processSkill(skill.id);
      await pollTask(task.id);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '操作失败';
      message.error(msg);
      setStatus('failed');
    } finally {
      setLoading(false);
    }
  };

  const statusInfo = status ? STATUS_MAP[status] ?? STATUS_MAP.pending : null;

  return (
    <div>
      <Title level={3}>
        <ExperimentOutlined /> 思维框架克隆
      </Title>
      <Paragraph type="secondary">
        上传你的文本语料，AI 将分析并生成属于你的思维 Skill 框架。
      </Paragraph>

      <Card>
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item
            label="框架名称"
            name="name"
            rules={[{ required: true, message: '请输入名称' }]}
          >
            <Input placeholder="例如：我的写作风格" maxLength={50} />
          </Form.Item>

          <Form.Item
            label="文本语料"
            name="corpus"
            rules={[{ required: true, message: '请输入或粘贴文本语料' }]}
          >
            <TextArea
              rows={10}
              placeholder="粘贴你的文章、对话记录、笔记等文本内容…"
              showCount
              maxLength={50000}
            />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                icon={<ExperimentOutlined />}
                size="large"
              >
                生成思维框架
              </Button>
              {statusInfo && (
                <Tag icon={statusInfo.icon} color={statusInfo.color}>
                  {statusInfo.text}
                </Tag>
              )}
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
