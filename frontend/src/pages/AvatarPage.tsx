import SkeletonCanvas from "../components/SkeletonCanvas";
import { useState } from 'react';
import {
  Button,
  Card,
  Input,
  Upload,
  Select,
  Tabs,
  Space,
  Typography,
  Image,
  Spin,
  Tag,
  message,
} from 'antd';
import {
  PictureOutlined,
  BoxPlotOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons';
import type { UploadFile } from 'antd/es/upload/interface';
import {
  uploadFile,
  createAvatar,
  processAvatar,
  getAvatar,
  getTask,
} from '../api/endpoints';
import type { Task } from '../api/types';
import type { PreviewState } from '../App';

const { Title, Paragraph, Text } = Typography;

const STYLE_OPTIONS = [
  { value: 'realistic', label: '写实风格' },
  { value: 'cartoon', label: '卡通风格' },
  { value: 'anime', label: '动漫风格' },
];

interface Props {
  setPreview: (p: PreviewState) => void;
}

export default function AvatarPage({ setPreview }: Props) {
  const [tab, setTab] = useState<'2d' | '3d'>('2d');
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [avatarName, setAvatarName] = useState('');
  const [style, setStyle] = useState('realistic');
  const [loading, setLoading] = useState(false);
  const [resultUrl, setResultUrl] = useState('');
  const [_animData, setAnimData] = useState<any>(null);
  const [status, setStatus] = useState('');
  const [_avatarId, setAvatarId] = useState('');

  const pollTask = async (taskId: string): Promise<boolean> => {
    let attempts = 0;
    while (attempts < 60) {
      await new Promise((r) => setTimeout(r, 2000));
      try {
        const { data: task } = await getTask(taskId) as { data: Task };
        if (task.status === 'done' || task.status === 'completed') return true;
        if (task.status === 'failed') {
          message.error(`生成失败: ${task.error || '未知错误'}`);
          return false;
        }
        setStatus(task.status);
      } catch {
        // continue
      }
      attempts++;
    }
    return false;
  };

  const handleGenerate = async () => {
    if (!fileList.length) {
      message.warning('请先上传照片');
      return;
    }
    if (!avatarName.trim()) {
      message.warning('请输入形象名称');
      return;
    }

    setLoading(true);
    setResultUrl('');
    setStatus('uploading');

    try {
      const rawFile = fileList[0].originFileObj as File;
      const { data: uploadRes } = await uploadFile(rawFile);

      setStatus('creating');
      const { data: avatar } = await createAvatar({
        name: avatarName,
        type: tab,
        photo_path: uploadRes.path,
        style: tab === '2d' ? style : undefined,
      });
      setAvatarId(avatar.id);

      setStatus('processing');
      const { data: task } = await processAvatar(avatar.id);
      const ok = await pollTask(task.id);

      if (ok) {
        const { data: a } = await getAvatar(avatar.id);
        setStatus('done');
        message.success('形象生成完成！');

        // Parse animation data from result
        let animData = null;
        let cartoonUrl = '';
        try {
          const parsed = JSON.parse(a.result);
          if (parsed.animation_data) {
            // Fetch animation JSON
            const animUrl = parsed.animation_data.replace(/^(.\/)?uploads\//, '/uploads/');
            const resp = await fetch(animUrl);
            animData = await resp.json();
          }
          if (parsed.cartoon_image) {
            cartoonUrl = parsed.cartoon_image.replace(/^(.\/)?uploads\//, '/uploads/');
          }
        } catch {
          // Result might be a plain URL
          cartoonUrl = a.result ? a.result.replace(/^(.\/)?uploads\//, '/uploads/') : '';
        }

        setResultUrl(cartoonUrl || '');
        setAnimData(animData);

        setPreview({
          visible: true,
          title: `🖼️ ${avatarName}`,
          content: (
            <div>
              {animData ? (
                <div>
                  <SkeletonCanvas
                    animationData={animData}
                    cartoonImageUrl={cartoonUrl}
                    width={400}
                    height={500}
                  />
                  <div style={{ marginTop: 12, color: '#888' }}>
                    骨骼动画已生成 · 15个关节 · 14根骨骼 · 可做动作
                  </div>
                </div>
              ) : cartoonUrl ? (
                <div>
                  <img src={cartoonUrl} alt="avatar" style={{ maxWidth: '100%', borderRadius: 8 }} />
                  <div style={{ marginTop: 12, color: '#888' }}>形象已生成</div>
                </div>
              ) : (
                <div style={{ textAlign: 'center', padding: 40, color: '#888' }}>
                  <PictureOutlined style={{ fontSize: 48, marginBottom: 16 }} />
                  <div>形象已创建</div>
                </div>
              )}
            </div>
          ),
        });
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '生成失败';
      message.error(msg);
      setStatus('failed');
    } finally {
      setLoading(false);
    }
  };

  const uploadProps = {
    fileList,
    beforeUpload: () => false,
    onChange: (info: { fileList: UploadFile[] }) => setFileList(info.fileList),
    maxCount: 1,
    accept: 'image/*',
  };

  const renderForm = () => (
    <Space direction="vertical" style={{ width: '100%' }} size="middle">
      <div>
        <Text strong>形象名称</Text>
        <Input
          placeholder="例如：我的虚拟形象"
          value={avatarName}
          onChange={(e) => setAvatarName(e.target.value)}
          maxLength={30}
          style={{ marginTop: 8 }}
        />
      </div>

      <div>
        <Text strong>上传照片</Text>
        <Upload.Dragger {...uploadProps} style={{ marginTop: 8 }}>
          <p className="ant-upload-drag-icon">
            <PictureOutlined />
          </p>
          <p className="ant-upload-text">点击或拖拽照片到此处上传</p>
          <p className="ant-upload-hint">建议使用正面清晰人脸照片</p>
        </Upload.Dragger>
      </div>

      {tab === '2d' && (
        <div>
          <Text strong>风格选择</Text>
          <Select
            value={style}
            onChange={setStyle}
            options={STYLE_OPTIONS}
            style={{ width: '100%', marginTop: 8 }}
          />
        </div>
      )}

      <Button
        type="primary"
        icon={tab === '2d' ? <PictureOutlined /> : <BoxPlotOutlined />}
        loading={loading}
        onClick={handleGenerate}
        size="large"
        block
      >
        {loading ? '生成中…' : `生成 ${tab === '2d' ? '2D' : '3D'} 形象`}
      </Button>

      {loading && (
        <div style={{ textAlign: 'center' }}>
          <Spin size="large" />
          <Paragraph style={{ marginTop: 8 }}>
            {status === 'uploading' && '上传照片中…'}
            {status === 'creating' && '创建形象记录…'}
            {status === 'processing' && 'AI 正在生成形象…'}
          </Paragraph>
        </div>
      )}

      {status === 'done' && (
        <Tag icon={<CheckCircleOutlined />} color="success">
          生成完成
        </Tag>
      )}
    </Space>
  );

  const renderResult = () => {
    if (!resultUrl) return null;
    const imgUrl = resultUrl.startsWith('/uploads/') ? resultUrl : resultUrl.replace(/^(\.\/)?uploads\//, '/uploads/');

    if (tab === '2d') {
      return (
        <Card title="🖼️ 生成结果" style={{ marginTop: 24 }}>
          <Image
            src={imgUrl}
            alt="2D Avatar"
            style={{ maxWidth: '100%', borderRadius: 8 }}
          />
        </Card>
      );
    }

    return (
      <Card title="🧊 3D 形象" style={{ marginTop: 24 }}>
        <div
          style={{
            height: 400,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: '#f0f2f5',
            borderRadius: 8,
          }}
        >
          <div style={{ textAlign: 'center' }}>
            <BoxPlotOutlined style={{ fontSize: 64, color: '#1677ff' }} />
            <Paragraph style={{ marginTop: 16 }}>
              3D 模型已生成，可在支持 WebGL 的浏览器中预览
            </Paragraph>
            <Button type="primary" href={resultUrl} target="_blank">
              打开 3D 预览
            </Button>
          </div>
        </div>
      </Card>
    );
  };

  const items = [
    {
      key: '2d',
      label: (
        <span>
          <PictureOutlined /> 2D 形象
        </span>
      ),
      children: (
        <>
          {renderForm()}
          {renderResult()}
        </>
      ),
    },
    {
      key: '3d',
      label: (
        <span>
          <BoxPlotOutlined /> 3D 形象
        </span>
      ),
      children: (
        <>
          {renderForm()}
          {renderResult()}
        </>
      ),
    },
  ];

  return (
    <div>
      <Title level={3}>
        <PictureOutlined /> 虚拟形象
      </Title>
      <Paragraph type="secondary">
        上传一张照片，生成你的 2D 或 3D 虚拟形象。
      </Paragraph>

      <Card>
        <Tabs
          activeKey={tab}
          onChange={(k) => {
            setTab(k as '2d' | '3d');
            setResultUrl('');
            setStatus('');
            setFileList([]);
          }}
          items={items}
        />
      </Card>
    </div>
  );
}
