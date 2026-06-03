import { useEffect, useState } from 'react';
import { Card, Row, Col, Tag, Typography, Spin, Empty, Button, message, Popconfirm } from 'antd';
import {
  ExperimentOutlined,
  SoundOutlined,
  PictureOutlined,
  DeleteOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { listSkills, listVoices, listAvatars, deleteSkill, deleteVoice, deleteAvatar } from '../api/endpoints';
import type { Skill, Voice, Avatar } from '../api/types';

const { Title, Text } = Typography;

const statusColor: Record<string, string> = {
  pending: 'default',
  processing: 'processing',
  training: 'processing',
  done: 'success',
  failed: 'error',
};

export default function Dashboard() {
  const [skills, setSkills] = useState<Skill[]>([]);
  const [voices, setVoices] = useState<Voice[]>([]);
  const [avatars, setAvatars] = useState<Avatar[]>([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    setLoading(true);
    try {
      const [s, v, a] = await Promise.all([
        listSkills(),
        listVoices(),
        listAvatars(),
      ]);
      setSkills(s.data || []);
      setVoices(v.data || []);
      setAvatars(a.data || []);
    } catch {
      message.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleDelete = async (type: string, id: string) => {
    try {
      if (type === 'skill') await deleteSkill(id);
      else if (type === 'voice') await deleteVoice(id);
      else if (type === 'avatar') await deleteAvatar(id);
      message.success('已删除');
      load();
    } catch {
      message.error('删除失败');
    }
  };

  if (loading) return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />;

  const total = skills.length + voices.length + avatars.length;

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>📊 我的数字分身</Title>
        <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
      </div>

      {total === 0 ? (
        <Empty description="还没有创建任何内容，去试试吧！" />
      ) : (
        <>
          {/* Skills */}
          {skills.length > 0 && (
            <>
              <Title level={4}><ExperimentOutlined /> 思维框架 ({skills.length})</Title>
              <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                {skills.map(s => (
                  <Col key={s.id} xs={24} sm={12} lg={8} xl={6}>
                    <Card
                      hoverable
                      title={s.name}
                      extra={<Tag color={statusColor[s.status]}>{s.status}</Tag>}
                      actions={[
                        <Popconfirm key="del" title="确认删除？" onConfirm={() => handleDelete('skill', s.id)}>
                          <DeleteOutlined />
                        </Popconfirm>,
                      ]}
                    >
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        {s.corpus?.slice(0, 80)}...
                      </Text>
                      <br />
                      <Text type="secondary" style={{ fontSize: 11 }}>
                        {new Date(s.created_at).toLocaleString('zh-CN')}
                      </Text>
                      {s.result && (
                        <div style={{ marginTop: 8 }}>
                          <Tag color="blue">已生成 SKILL.md</Tag>
                        </div>
                      )}
                    </Card>
                  </Col>
                ))}
              </Row>
            </>
          )}

          {/* Voices */}
          {voices.length > 0 && (
            <>
              <Title level={4}><SoundOutlined /> 语音克隆 ({voices.length})</Title>
              <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                {voices.map(v => (
                  <Col key={v.id} xs={24} sm={12} lg={8} xl={6}>
                    <Card
                      hoverable
                      title={v.name}
                      extra={<Tag color={statusColor[v.status]}>{v.status}</Tag>}
                      actions={[
                        <Popconfirm key="del" title="确认删除？" onConfirm={() => handleDelete('voice', v.id)}>
                          <DeleteOutlined />
                        </Popconfirm>,
                      ]}
                    >
                      <Text type="secondary">时长: {v.duration?.toFixed(1) || 0}s</Text>
                      <br />
                      <Text type="secondary" style={{ fontSize: 11 }}>
                        {new Date(v.created_at).toLocaleString('zh-CN')}
                      </Text>
                      {v.status === 'done' && (
                        <div style={{ marginTop: 8 }}>
                          <Tag color="green">模型已就绪</Tag>
                        </div>
                      )}
                    </Card>
                  </Col>
                ))}
              </Row>
            </>
          )}

          {/* Avatars */}
          {avatars.length > 0 && (
            <>
              <Title level={4}><PictureOutlined /> 虚拟形象 ({avatars.length})</Title>
              <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                {avatars.map(a => (
                  <Col key={a.id} xs={24} sm={12} lg={8} xl={6}>
                    <Card
                      hoverable
                      title={a.name}
                      extra={
                        <>
                          <Tag>{a.type === '2d' ? '2D' : '3D'}</Tag>
                          <Tag color={statusColor[a.status]}>{a.status}</Tag>
                        </>
                      }
                      cover={
                        a.result ? (
                          <img
                            alt={a.name}
                            src={a.result.replace('./uploads/', '/uploads/')}
                            style={{ height: 160, objectFit: 'cover' }}
                          />
                        ) : undefined
                      }
                      actions={[
                        <Popconfirm key="del" title="确认删除？" onConfirm={() => handleDelete('avatar', a.id)}>
                          <DeleteOutlined />
                        </Popconfirm>,
                      ]}
                    >
                      <Text type="secondary">风格: {a.style || '写实'}</Text>
                      <br />
                      <Text type="secondary" style={{ fontSize: 11 }}>
                        {new Date(a.created_at).toLocaleString('zh-CN')}
                      </Text>
                      {a.status === 'done' && (
                        <div style={{ marginTop: 8 }}>
                          <Tag color="green">已生成</Tag>
                        </div>
                      )}
                    </Card>
                  </Col>
                ))}
              </Row>
            </>
          )}
        </>
      )}
    </div>
  );
}
