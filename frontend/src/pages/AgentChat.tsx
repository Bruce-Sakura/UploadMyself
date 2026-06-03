import { useState, useRef, useEffect } from 'react';
import { Input, Button, Typography, Spin, Tag } from 'antd';
import { SendOutlined, RobotOutlined, UserOutlined } from '@ant-design/icons';

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
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

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

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: 'calc(100vh - 160px)' }}>
      <div style={{ marginBottom: 16 }}>
        <Text strong style={{ fontSize: 18 }}>🤖 对话</Text>
        <Text type="secondary" style={{ marginLeft: 12 }}>
          和你的数字分身聊天（需要配置 LLM_API_KEY）
        </Text>
      </div>

      {/* 消息区 */}
      <div
        style={{
          flex: 1,
          overflow: 'auto',
          border: '1px solid #f0f0f0',
          borderRadius: 8,
          padding: 16,
          background: '#fafafa',
        }}
      >
        {messages.length === 0 && (
          <div style={{ textAlign: 'center', color: '#bbb', marginTop: 80 }}>
            <RobotOutlined style={{ fontSize: 48, marginBottom: 16 }} />
            <div>开始和你的数字分身对话吧</div>
            <div style={{ fontSize: 12, marginTop: 8 }}>
              需要先在 docker-compose.yml 中配置 LLM_API_KEY
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
            }}
          >
            <div
              style={{
                maxWidth: '70%',
                padding: '10px 14px',
                borderRadius: 12,
                background: msg.role === 'user' ? '#1890ff' : '#fff',
                color: msg.role === 'user' ? '#fff' : '#333',
                border: msg.role === 'user' ? 'none' : '1px solid #e8e8e8',
                boxShadow: '0 1px 2px rgba(0,0,0,0.06)',
              }}
            >
              <div style={{ fontSize: 12, marginBottom: 4, opacity: 0.7 }}>
                {msg.role === 'user' ? <UserOutlined /> : <RobotOutlined />}
              </div>
              <div style={{ whiteSpace: 'pre-wrap' }}>{msg.content}</div>
              {msg.toolCalls && msg.toolCalls.length > 0 && (
                <div style={{ marginTop: 8, borderTop: '1px solid #f0f0f0', paddingTop: 8 }}>
                  {msg.toolCalls.map((tc, j) => (
                    <Tag key={j} color="blue" style={{ marginBottom: 4 }}>
                      🔧 {tc.content.slice(0, 60)}
                    </Tag>
                  ))}
                </div>
              )}
            </div>
          </div>
        ))}

        {loading && (
          <div style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: 12 }}>
            <div style={{ padding: '10px 14px', borderRadius: 12, background: '#fff', border: '1px solid #e8e8e8' }}>
              <Spin size="small" /> <Text type="secondary" style={{ marginLeft: 8 }}>思考中...</Text>
            </div>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* 输入区 */}
      <div style={{ display: 'flex', gap: 8, marginTop: 12 }}>
        <Input.TextArea
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder="输入消息... (Enter 发送, Shift+Enter 换行)"
          autoSize={{ minRows: 1, maxRows: 4 }}
          onPressEnter={e => {
            if (!e.shiftKey) {
              e.preventDefault();
              sendMessage();
            }
          }}
          disabled={loading}
        />
        <Button
          type="primary"
          icon={<SendOutlined />}
          onClick={sendMessage}
          loading={loading}
          style={{ height: 'auto' }}
        >
          发送
        </Button>
      </div>
    </div>
  );
}
