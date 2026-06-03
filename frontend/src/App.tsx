import { Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Typography } from 'antd';
import {
  HomeOutlined,
  ExperimentOutlined,
  SoundOutlined,
  PictureOutlined,
  DashboardOutlined,
  RobotOutlined,
} from '@ant-design/icons';
import SkillClone from './pages/SkillClone';
import VoiceClone from './pages/VoiceClone';
import AvatarPage from './pages/AvatarPage';
import Dashboard from './pages/Dashboard';
import AgentChat from './pages/AgentChat';
import { useState } from 'react';

const { Header, Content, Sider } = Layout;
const { Title } = Typography;

const menuItems = [
  { key: '/', icon: <HomeOutlined />, label: '首页' },
  { key: '/dashboard', icon: <DashboardOutlined />, label: '我的分身' },
  { key: '/chat', icon: <RobotOutlined />, label: '对话' },
  { key: '/skill', icon: <ExperimentOutlined />, label: '思维克隆' },
  { key: '/voice', icon: <SoundOutlined />, label: '语音克隆' },
  { key: '/avatar', icon: <PictureOutlined />, label: '虚拟形象' },
];

// 右侧预览面板的全局状态
export interface PreviewState {
  visible: boolean;
  title: string;
  content: React.ReactNode;
}

function App() {
  const navigate = useNavigate();
  const location = useLocation();
  const [preview, setPreview] = useState<PreviewState>({
    visible: false,
    title: '',
    content: null,
  });

  // 是否显示右侧预览（skill/voice/avatar 页面才显示）
  const showPreview = ['/skill', '/voice', '/avatar'].includes(location.pathname);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" width={220} breakpoint="lg" collapsedWidth={80}>
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: 18,
            fontWeight: 'bold',
            letterSpacing: 1,
          }}
        >
          🧬 UploadMyself
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          <Title level={4} style={{ margin: 0 }}>
            克隆你自己
          </Title>
        </Header>
        <Content style={{ margin: 0, display: 'flex', overflow: 'hidden' }}>
          {/* 左侧：操作区 */}
          <div
            style={{
              flex: showPreview ? '1 1 50%' : '1 1 100%',
              padding: 24,
              overflow: 'auto',
              background: '#fff',
            }}
          >
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/chat" element={<AgentChat />} />
              <Route path="/skill" element={<SkillClone setPreview={setPreview} />} />
              <Route path="/voice" element={<VoiceClone setPreview={setPreview} />} />
              <Route path="/avatar" element={<AvatarPage setPreview={setPreview} />} />
            </Routes>
          </div>

          {/* 右侧：预览/展示区 */}
          {showPreview && (
            <div
              style={{
                flex: '0 0 50%',
                borderLeft: '1px solid #f0f0f0',
                background: '#fafafa',
                padding: 24,
                overflow: 'auto',
              }}
            >
              {preview.visible ? (
                <div>
                  <Title level={4} style={{ marginTop: 0, marginBottom: 16 }}>
                    {preview.title}
                  </Title>
                  {preview.content}
                </div>
              ) : (
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    height: '100%',
                    color: '#bbb',
                    fontSize: 16,
                  }}
                >
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ fontSize: 48, marginBottom: 16 }}>👁️</div>
                    <div>在左侧操作，结果会在这里显示</div>
                  </div>
                </div>
              )}
            </div>
          )}
        </Content>
      </Layout>
    </Layout>
  );
}

function HomePage() {
  const navigate = useNavigate();
  return (
    <div>
      <Title level={2}>🧬 UploadMyself</Title>
      <Title level={4} type="secondary">
        上传你的数据，生成数字分身
      </Title>
      <div style={{ marginTop: 32 }}>
        {[
          {
            icon: '🤖',
            title: '对话',
            desc: '和你的数字分身聊天，它用你的方式思考',
            path: '/chat',
          },
          {
            icon: '📊',
            title: '我的分身',
            desc: '查看所有已创建的思维/声音/形象',
            path: '/dashboard',
          },
          {
            icon: '🧠',
            title: '思维框架克隆',
            desc: '文本语料 → 你的思维方式',
            path: '/skill',
          },
          {
            icon: '🎤',
            title: '语音克隆',
            desc: '语音样本 → 你的声音',
            path: '/voice',
          },
          {
            icon: '🖼️',
            title: '虚拟形象',
            desc: '一张照片 → 你的数字外貌',
            path: '/avatar',
          },
        ].map((item) => (
          <div
            key={item.title}
            onClick={() => navigate(item.path)}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 16,
              padding: '16px',
              marginBottom: 8,
              borderRadius: 8,
              border: '1px solid #f0f0f0',
              cursor: 'pointer',
              transition: 'all 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.borderColor = '#1890ff';
              e.currentTarget.style.background = '#f6ffed';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.borderColor = '#f0f0f0';
              e.currentTarget.style.background = 'transparent';
            }}
          >
            <span style={{ fontSize: 32 }}>{item.icon}</span>
            <div>
              <div style={{ fontWeight: 600, fontSize: 16 }}>{item.title}</div>
              <div style={{ color: '#888' }}>{item.desc}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
