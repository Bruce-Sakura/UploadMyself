import { useState } from 'react';
import {
  Button,
  Card,
  Input,
  Form,
  Tag,
  Spin,
  Upload,
  Space,
  Typography,
  Divider,
  message,
} from 'antd';
import {
  ExperimentOutlined,
  CheckCircleOutlined,
  LoadingOutlined,
  CloseCircleOutlined,
  FileTextOutlined,
  UploadOutlined,
  FilePdfOutlined,
  FileImageOutlined,
} from '@ant-design/icons';
import {
  createSkill,
  processSkill,
  getSkill,
  uploadCorpus,
} from '../api/endpoints';
import { pollTask } from '../hooks/useTaskPoller';
import type { UploadFile } from 'antd/es/upload/interface';
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
  const [status, setStatus] = useState('');
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [extracting, setExtracting] = useState(false);
  const [corpus, setCorpus] = useState('');

  const showInPreview = (title: string, content: React.ReactNode) => {
    setPreview({ visible: true, title, content });
  };

  const handlePoll = async (taskId: string) => {
    await pollTask(taskId, {
      onProgress: (task) => {
        setStatus(task.status);
        showInPreview('⏳ 处理中...', (
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin size="large" />
            <div style={{ marginTop: 16, color: '#888' }}>
              AI 正在分析你的语料，生成思维框架…
            </div>
          </div>
        ));
      },
      onDone: async (task) => {
        const { data: skill } = await getSkill(task.ref_id);
        setStatus('done');
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
      },
      onFailed: (task) => {
        setStatus('failed');
        message.error(`处理失败: ${task.error || '未知错误'}`);
      },
    });
  };

  const handleFileExtract = async () => {
    if (!fileList.length) {
      message.warning('请先选择文件');
      return;
    }
    const rawFile = fileList[0].originFileObj as File;
    setExtracting(true);
    try {
      const { data } = await uploadCorpus(rawFile);
      const newText = data.text || '';
      setCorpus((prev) => {
        const combined = prev ? prev + '\n\n' + newText : newText;
        form.setFieldsValue({ corpus: combined });
        return combined;
      });
      const methodLabel: Record<string, string> = {
        pdf: 'PDF 文本提取',
        docx: 'Word 文档提取',
        ocr: 'OCR 文字识别',
        paddleocr: 'PaddleOCR 识别',
        tesseract: 'Tesseract OCR',
        text: '纯文本读取',
      };
      message.success(`${methodLabel[data.method] || '文本提取'}完成！已提取 ${newText.length} 字`);
      setFileList([]);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '文件提取失败';
      message.error(msg);
    } finally {
      setExtracting(false);
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
      await handlePoll(task.id);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '操作失败';
      message.error(msg);
      setStatus('failed');
    } finally {
      setLoading(false);
    }
  };

  const statusInfo = status ? STATUS_MAP[status] ?? STATUS_MAP.pending : null;

  const uploadProps = {
    fileList,
    beforeUpload: () => false,
    onChange: (info: { fileList: UploadFile[] }) => setFileList(info.fileList),
    maxCount: 1,
    accept: '.pdf,.doc,.docx,.txt,.md,.png,.jpg,.jpeg,.webp',
  };

  const getFileInfo = () => {
    if (!fileList.length) return null;
    const f = fileList[0];
    const name = f.name || '';
    const ext = name.split('.').pop()?.toLowerCase() || '';
    if (ext === 'pdf') return { icon: <FilePdfOutlined />, color: '#ff4d4f', label: 'PDF 文档' };
    if (['doc', 'docx'].includes(ext)) return { icon: <FileTextOutlined />, color: '#1677ff', label: 'Word 文档' };
    if (['png', 'jpg', 'jpeg', 'webp'].includes(ext)) return { icon: <FileImageOutlined />, color: '#52c41a', label: '图片 (OCR)' };
    return { icon: <FileTextOutlined />, color: '#888', label: '文本文件' };
  };

  const fileInfo = getFileInfo();

  return (
    <div>
      <Title level={3}>
        <ExperimentOutlined /> 思维框架克隆
      </Title>
      <Paragraph type="secondary">
        上传你的文本语料或文件（PDF/Word/图片），AI 将分析并生成属于你的思维 Skill 框架。
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

          {/* File Upload Section */}
          <Form.Item label="从文件导入语料">
            <Space.Compact style={{ width: '100%' }}>
              <Upload.Dragger
                {...uploadProps}
                style={{ padding: '12px 16px', flex: 1 }}
              >
                <p style={{ margin: 0 }}>
                  <UploadOutlined style={{ marginRight: 8 }} />
                  点击或拖拽文件到此处
                </p>
                <p style={{ margin: '4px 0 0', color: '#999', fontSize: 12 }}>
                  支持 PDF、Word (.docx)、图片 (OCR)、纯文本
                </p>
              </Upload.Dragger>
              <Button
                type="primary"
                icon={extracting ? <LoadingOutlined /> : <ExperimentOutlined />}
                loading={extracting}
                onClick={handleFileExtract}
                disabled={!fileList.length}
                style={{ height: 'auto', minHeight: 80 }}
              >
                {extracting ? '提取中…' : '提取文本'}
              </Button>
            </Space.Compact>

            {fileInfo && (
              <div style={{ marginTop: 8 }}>
                <Tag icon={fileInfo.icon} color={fileInfo.color}>
                  {fileInfo.label}: {fileList[0]?.name}
                </Tag>
              </div>
            )}
          </Form.Item>

          <Divider plain>或直接粘贴文本</Divider>

          <Form.Item
            label="文本语料"
            name="corpus"
            rules={[{ required: true, message: '请输入或粘贴文本语料' }]}
          >
            <TextArea
              rows={10}
              placeholder="粘贴你的文章、对话记录、笔记等文本内容…&#10;&#10;也可以通过上方文件上传自动提取文本"
              showCount
              maxLength={50000}
              value={corpus}
              onChange={(e) => setCorpus(e.target.value)}
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
