import { useState, useRef } from 'react';
import {
  Button,
  Card,
  Input,
  Upload,
  Progress,
  Space,
  Typography,
  message,
  Divider,
} from 'antd';
import {
  SoundOutlined,
  AudioOutlined,
  InboxOutlined,
} from '@ant-design/icons';
import type { UploadFile } from 'antd/es/upload/interface';
import { uploadFile, createVoice, trainVoice, synthesizeVoice, getTask, getVoice } from '../api/endpoints';
import type { Task } from '../api/types';

const { Title, Paragraph, Text } = Typography;
const { Dragger } = Upload;
const { TextArea } = Input;

export default function VoiceClone() {
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [voiceName, setVoiceName] = useState('');
  const [training, setTraining] = useState(false);
  const [progress, setProgress] = useState(0);
  const [voiceId, setVoiceId] = useState('');
  const [voiceStatus, setVoiceStatus] = useState('');

  // Synthesize section
  const [synthText, setSynthText] = useState('');
  const [synthesizing, setSynthesizing] = useState(false);
  const [audioUrl, setAudioUrl] = useState('');
  const audioRef = useRef<HTMLAudioElement>(null);

  const pollTask = async (taskId: string): Promise<boolean> => {
    let attempts = 0;
    while (attempts < 60) {
      await new Promise((r) => setTimeout(r, 2000));
      try {
        const { data: task } = await getTask(taskId) as { data: Task };
        setProgress(task.progress ?? 0);
        if (task.status === 'done' || task.status === 'completed') {
          return true;
        }
        if (task.status === 'failed') {
          message.error(`训练失败: ${task.error || '未知错误'}`);
          return false;
        }
      } catch {
        // continue
      }
      attempts++;
    }
    return false;
  };

  const handleTrain = async () => {
    if (!fileList.length) {
      message.warning('请先上传音频文件');
      return;
    }
    if (!voiceName.trim()) {
      message.warning('请输入声音名称');
      return;
    }

    setTraining(true);
    setProgress(0);
    setVoiceStatus('uploading');

    try {
      // 1. Upload audio
      const rawFile = fileList[0].originFileObj as File;
      const { data: uploadRes } = await uploadFile(rawFile);

      // 2. Create voice
      setVoiceStatus('creating');
      const { data: voice } = await createVoice({
        name: voiceName,
        audio_path: uploadRes.path,
      });
      setVoiceId(voice.id);

      // 3. Train
      setVoiceStatus('training');
      const { data: task } = await trainVoice(voice.id);
      const ok = await pollTask(task.id);

      if (ok) {
        const { data: v } = await getVoice(voice.id);
        setVoiceStatus(v.status);
        message.success('声音训练完成！');
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '训练失败';
      message.error(msg);
      setVoiceStatus('failed');
    } finally {
      setTraining(false);
    }
  };

  const handleSynthesize = async () => {
    if (!voiceId) {
      message.warning('请先完成声音训练');
      return;
    }
    if (!synthText.trim()) {
      message.warning('请输入要合成的文本');
      return;
    }

    setSynthesizing(true);
    setAudioUrl('');
    try {
      const { data } = await synthesizeVoice(voiceId, synthText);
      setAudioUrl(data.audio_url);
      message.success('语音合成完成！');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '合成失败';
      message.error(msg);
    } finally {
      setSynthesizing(false);
    }
  };

  const uploadProps = {
    fileList,
    beforeUpload: () => false,
    onChange: (info: { fileList: UploadFile[] }) => setFileList(info.fileList),
    maxCount: 1,
    accept: 'audio/*',
  };

  return (
    <div>
      <Title level={3}>
        <SoundOutlined /> 语音克隆
      </Title>
      <Paragraph type="secondary">
        上传一段语音样本，训练属于你的克隆声音，然后用它合成任意文本。
      </Paragraph>

      {/* Training Section */}
      <Card title="🎤 训练声音" style={{ marginBottom: 24 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <div>
            <Text strong>声音名称</Text>
            <Input
              placeholder="例如：我的声音"
              value={voiceName}
              onChange={(e) => setVoiceName(e.target.value)}
              maxLength={30}
              style={{ marginTop: 8 }}
            />
          </div>

          <div>
            <Text strong>上传音频</Text>
            <Dragger {...uploadProps} style={{ marginTop: 8 }}>
              <p className="ant-upload-drag-icon">
                <InboxOutlined />
              </p>
              <p className="ant-upload-text">点击或拖拽音频文件到此处上传</p>
              <p className="ant-upload-hint">
                支持 WAV、MP3、FLAC 等格式，建议 10 秒以上清晰录音
              </p>
            </Dragger>
          </div>

          <Button
            type="primary"
            icon={<AudioOutlined />}
            loading={training}
            onClick={handleTrain}
            size="large"
            block
          >
            {training ? '训练中…' : '开始训练'}
          </Button>

          {training && (
            <div>
              <Text type="secondary">
                状态：{voiceStatus === 'uploading' && '上传音频中…'}
                {voiceStatus === 'creating' && '创建声音记录…'}
                {voiceStatus === 'training' && '训练模型中…'}
              </Text>
              <Progress percent={progress} status="active" />
            </div>
          )}

          {voiceId && !training && voiceStatus !== 'failed' && (
            <Paragraph type="success">
              ✅ 声音训练完成 — ID: {voiceId}
            </Paragraph>
          )}
        </Space>
      </Card>

      {/* Synthesize Section */}
      <Card title="🗣️ 语音合成">
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <TextArea
            rows={4}
            placeholder="输入要合成的文本…"
            value={synthText}
            onChange={(e) => setSynthText(e.target.value)}
            showCount
            maxLength={5000}
          />

          <Button
            type="primary"
            icon={<SoundOutlined />}
            loading={synthesizing}
            onClick={handleSynthesize}
            disabled={!voiceId}
            size="large"
          >
            {synthesizing ? '合成中…' : '开始合成'}
          </Button>

          {audioUrl && (
            <>
              <Divider />
              <div>
                <Text strong>合成结果</Text>
                <audio
                  ref={audioRef}
                  controls
                  src={audioUrl}
                  style={{ width: '100%', marginTop: 8 }}
                />
              </div>
            </>
          )}
        </Space>
      </Card>
    </div>
  );
}
