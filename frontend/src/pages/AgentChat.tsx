import { useState, useRef, useEffect } from 'react';
import { Input, Button, Typography, Spin, Tag, Select, Card, Avatar as AntAvatar } from 'antd';
import { SendOutlined, RobotOutlined, UserOutlined, PictureOutlined } from '@ant-design/icons';
import { listSkills, listAvatars } from '../api/endpoints';
import type { Skill, Avatar } from '../api/types';

const { Text } = Typography;

interface ChatMsg {
  role: 'user' | 'assistant';
  content: string;
  toolCalls?: { tool_call_id: string; content: string }[];
  timestamp: string;
}

export default function AgentChat() {
  const [messages, setMessages] = useState<ChatMsg[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [conversationId, setConversationId] = useState('');

  // 分身选择
  const [skills, setSkills] = useState<Skill[]>([]);
  const [avatars, setAvatars] = useState<Avatar[]>([]);
  const [selectedSkill, setSelectedSkill] = useState<string>('');
  const [selectedAvatar, setSelectedAvatar] = useState<Avatar | null>(null);

  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // 加载可用的 skill 和 avatar
  useEffect(() => {
    listSkills().then(res => {
      const doneSkills = (res.data || []).filter((s: Skill) => s.status === 'done' && s.result);
      setSkills(doneSkills);
    }).catch(() => {});
    listAvatars().then(res => {
      const doneAvatars = (res.data || []).filter((a: Avatar) => a.status === 'done');
      setAvatars(doneAvatars);
    }).catch(() => {});
  }, []);

  const sendMessage = async () => {
    if (!input.trim() || loading) return;

    const userMsg: ChatMsg = { role: 'user', content: input, timestamp: new Date().toISOString() };
    setMessages(prev => [...prev, userMsg]);
    setInput('');
    setLoading(true);

    try {
      const resp = await fetch('/api/v1/agent/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message: userMsg.content,
          conversation_id: conversationId,
          skill_id: selectedSkill || undefined,
        }),
      });

      const data = await resp.json();

      if (data.error) {
        setMessages(prev => [...prev, {
          role: 'assistant',
          content: `❌ 错误: ${data.error}`,
          timestamp: new Date().toISOString(),
        }]);
      } else {
        if (data.conversation_id) setConversationId(data.conversation_id);
        setMessages(prev => [...prev, {
          role: 'assistant',
          content: data.reply,
          toolCalls: data.tool_calls,
          timestamp: data.timestamp,
        }]);
      }
    } catch (err: any) {
      setMessages(prev => [...prev, {
        role: 'assistant',
        content: `❌ 网络错误: ${err.message}`,
        timestamp: new Date().toISOString(),
      }]);
    } finally {
      setLoading(false);
    }
  };

  // 获取选中 avatar 的图片 URL
  const avatarUrl = selectedAvatar?.result
    ? selectedAvatar.result.replace(/^(.\/)?uploads\//, '/uploads/')
    : '';

  return (
    <div style={{ display: 'flex', gap: 16, height: 'calc(100vh - 160px)' }}>
      {/* 左侧：选择面板 */}
      <div style={{ width: 240, flexShrink: 0 }}>
        <Card size="small" title="🧠 选择思维" style={{ marginBottom: 12 }}>
          {skills.length === 0 ? (
            <Text type="secondary" style={{ fontSize: 12 }}>
              还没有已生成的思维框架，去「思维克隆」创建一个
            </Text>
          ) : (
            <Select
              placeholder="选择一个思维框架"
              style={{ width: '100%' }}
              value={selectedSkill || undefined}
              onChange={setSelectedSkill}
              allowClear
              options={skills.map(s => ({ value: s.id, label: s.name }))}
            />
          )}
          {selectedSkill && (
            <Tag color="blue" style={{ marginTop: 8 }}>
              ✅ 已加载思维框架
            </Tag>
          )}
        </Card>

        <Card size="small" title="🖼️ 选择形象" style={{ marginBottom: 12 }}>
          {avatars.length === 0 ? (
            <Text type="secondary" style={{ fontSize: 12 }}>
              还没有已生成的形象，去「虚拟形象」创建一个
            </Text>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {avatars.map(a => {
                const url = a.result ? a.result.replace(/^(.\/)?uploads\//, '/uploads/') : '';
                const isSelected = selectedAvatar?.id === a.id;
                return (
                  <div
                    key={a.id}
                    onClick={() => setSelectedAvatar(isSelected ? null : a)}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 8,
                      padding: 8,
                      borderRadius: 8,
                      border: isSelected ? '2px solid #1890ff' : '1px solid #f0f0f0',
                      cursor: 'pointer',
                      background: isSelected ? '#e6f7ff' : 'transparent',
                    }}
                  >
                    {url ? (
                      <img src={url} alt={a.name} style={{ width: 40, height: 40, borderRadius: '50%', objectFit: 'cover' }} />
                    ) : (
                      <AntAvatar icon={<PictureOutlined />} size={40} />
                    )}
                    <div>
                      <div style={{ fontWeight: 500, fontSize: 13 }}>{a.name}</div>
                      <Tag>{a.type === '2d' ? '2D' : '3D'}</Tag>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </Card>

        <Card size="small" title="📊 对话信息">
          <div style={{ fontSize: 12, color: '#888' }}>
            <div>消息数: {messages.length}</div>
            <div>会话ID: {conversationId ? conversationId.slice(0, 8) + '...' : '未开始'}</div>
            <div>思维框架: {selectedSkill ? '已加载' : '默认'}</div>
            <div>形象: {selectedAvatar ? selectedAvatar.name : '未选择'}</div>
          </div>
        </Card>
      </div>

      {/* 右侧：对话区 */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        {/* 对话标题栏 */}
        <div style={{
          padding: '8px 16px',
          background: '#fafafa',
          borderRadius: '8px 8px 0 0',
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          gap: 12,
        }}>
          {avatarUrl ? (
            <img src={avatarUrl} alt="avatar" style={{ width: 36, height: 36, borderRadius: '50%', objectFit: 'cover' }} />
          ) : (
            <AntAvatar icon={<RobotOutlined />} size={36} />
          )}
          <div>
            <Text strong>{selectedAvatar?.name || '数字分身'}</Text>
            <br />
            <Text type="secondary" style={{ fontSize: 11 }}>
              {selectedSkill ? '使用自定义思维框架' : '使用默认思维'}
            </Text>
          </div>
        </div>

        {/* 消息区 */}
        <div style={{
          flex: 1,
          overflow: 'auto',
          borderLeft: '1px solid #f0f0f0',
          borderRight: '1px solid #f0f0f0',
          padding: 16,
          background: '#fff',
        }}>
          {messages.length === 0 && (
            <div style={{ textAlign: 'center', color: '#bbb', marginTop: 80 }}>
              {avatarUrl ? (
                <img src={avatarUrl} alt="avatar" style={{ width: 80, height: 80, borderRadius: '50%', objectFit: 'cover', marginBottom: 16 }} />
              ) : (
                <RobotOutlined style={{ fontSize: 48, marginBottom: 16 }} />
              )}
              <div style={{ fontSize: 16, marginBottom: 8 }}>
                {selectedAvatar ? `和 ${selectedAvatar.name} 对话` : '开始对话'}
              </div>
              <div style={{ fontSize: 12 }}>
                {selectedSkill ? '已加载自定义思维框架' : '左侧选择思维框架和形象'}
              </div>
            </div>
          )}

          {messages.map((msg, i) => (
            <div
              key={i}
              style={{
                display: 'flex',
                justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
                marginBottom: 12,
                gap: 8,
              }}
            >
              {msg.role === 'assistant' && (
                avatarUrl ? (
                  <img src={avatarUrl} alt="" style={{ width: 32, height: 32, borderRadius: '50%', objectFit: 'cover', flexShrink: 0 }} />
                ) : (
                  <AntAvatar icon={<RobotOutlined />} size={32} style={{ flexShrink: 0 }} />
                )
              )}
              <div
                style={{
                  maxWidth: '70%',
                  padding: '10px 14px',
                  borderRadius: 12,
                  background: msg.role === 'user' ? '#1890ff' : '#f6f6f6',
                  color: msg.role === 'user' ? '#fff' : '#333',
                  boxShadow: '0 1px 2px rgba(0,0,0,0.06)',
                }}
              >
                <div style={{ whiteSpace: 'pre-wrap' }}>{msg.content}</div>
                {msg.toolCalls && msg.toolCalls.length > 0 && (
                  <div style={{ marginTop: 8, borderTop: '1px solid #e8e8e8', paddingTop: 8 }}>
                    {msg.toolCalls.map((tc, j) => (
                      <Tag key={j} color="blue" style={{ marginBottom: 4 }}>🔧 {tc.content.slice(0, 60)}</Tag>
                    ))}
                  </div>
                )}
              </div>
              {msg.role === 'user' && (
                <AntAvatar icon={<UserOutlined />} size={32} style={{ flexShrink: 0 }} />
              )}
            </div>
          ))}

          {loading && (
            <div style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: 12, gap: 8 }}>
              {avatarUrl ? (
                <img src={avatarUrl} alt="" style={{ width: 32, height: 32, borderRadius: '50%', objectFit: 'cover' }} />
              ) : (
                <AntAvatar icon={<RobotOutlined />} size={32} />
              )}
              <div style={{ padding: '10px 14px', borderRadius: 12, background: '#f6f6f6' }}>
                <Spin size="small" /> <Text type="secondary" style={{ marginLeft: 8 }}>思考中...</Text>
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        {/* 输入区 */}
        <div style={{
          display: 'flex',
          gap: 8,
          padding: 12,
          border: '1px solid #f0f0f0',
          borderTop: 'none',
          borderRadius: '0 0 8px 8px',
          background: '#fafafa',
        }}>
          <Input.TextArea
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder={selectedAvatar ? `和 ${selectedAvatar.name} 说点什么...` : '输入消息...'}
            autoSize={{ minRows: 1, maxRows: 4 }}
            onPressEnter={e => {
              if (!e.shiftKey) { e.preventDefault(); sendMessage(); }
            }}
            disabled={loading}
          />
          <Button type="primary" icon={<SendOutlined />} onClick={sendMessage} loading={loading} style={{ height: 'auto' }}>
            发送
          </Button>
        </div>
      </div>
    </div>
  );
}
