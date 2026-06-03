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
  CopyOutlined,
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

const { TextArea } = Input;
const { Title, Paragraph } = Typography;

const STATUS_MAP: Record<string, { color: string; icon: React.ReactNode; text: string }> = {
  pending: { color: 'default', icon: <LoadingOutlined />, text: '等待中' },
  processing: { color: 'processing', icon: <LoadingOutlined />, text: '处理中' },
  done: { color: 'success', icon: <CheckCircleOutlined />, text: '已完成' },
  completed: { color: 'success', icon: <CheckCircleOutlined />, text: '已完成' },
  failed: { color: 'error', icon: <CloseCircleOutlined />, text: '失败' },
};

export default function SkillClone() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string>('');
  const [status, setStatus] = useState<string>('');
  const [skillId, setSkillId] = useState<string>('');

  const pollTask = async (taskId: string) => {
    let attempts = 0;
    const maxAttempts = 60;
    while (attempts < maxAttempts) {
      await new Promise((r) => setTimeout(r, 2000));
      try {
        const { data: task } = await getTask(taskId) as { data: Task };
        if (task.status === 'done' || task.status === 'completed') {
          const { data: skill } = await getSkill(task.ref_id);
          setResult(skill.result);
          setStatus('done');
          message.success('思维框架生成成功！');
          return;
        }
        if (task.status === 'failed') {
          setStatus('failed');
          message.error(`处理失败: ${task.error || '未知错误'}`);
          return;
        }
        setStatus(task.status);
      } catch {
        // continue polling
      }
      attempts++;
    }
  };

  const onFinish = async (values: { name: string; corpus: string }) => {
    setLoading(true);
    setResult('');
    setStatus('pending');
    try {
      const { data: skill } = await createSkill(values);
      setSkillId(skill.id);
      setStatus('processing');

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

  const copyResult = () => {
    navigator.clipboard.writeText(result);
    message.success('已复制到剪贴板');
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

      <Card style={{ marginBottom: 24 }}>
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

      {loading && status === 'processing' && (
        <Card>
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin size="large" />
            <Paragraph style={{ marginTop: 16 }}>
              AI 正在分析你的语料，生成思维框架…
            </Paragraph>
            {skillId && (
              <Paragraph type="secondary" style={{ fontSize: 12 }}>
                Skill ID: {skillId}
              </Paragraph>
            )}
          </div>
        </Card>
      )}

      {result && (
        <Card
          title="生成结果 — SKILL.md"
          extra={
            <Button icon={<CopyOutlined />} onClick={copyResult}>
              复制
            </Button>
          }
        >
          <pre
            style={{
              background: '#f6f8fa',
              padding: 16,
              borderRadius: 8,
              overflow: 'auto',
              maxHeight: 500,
              fontSize: 13,
              lineHeight: 1.6,
            }}
          >
            {result}
          </pre>
        </Card>
      )}
    </div>
  );
}
