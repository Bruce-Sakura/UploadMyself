import { Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Typography } from 'antd';
import {
  HomeOutlined,
  ExperimentOutlined,
  SoundOutlined,
  PictureOutlined,
  DashboardOutlined,
} from '@ant-design/icons';
import SkillClone from './pages/SkillClone';
import VoiceClone from './pages/VoiceClone';
import AvatarPage from './pages/AvatarPage';
import Dashboard from './pages/Dashboard';

const { Header, Content, Sider } = Layout;
const { Title } = Typography;

const menuItems = [
  { key: '/', icon: <HomeOutlined />, label: '首页' },
  { key: '/dashboard', icon: <DashboardOutlined />, label: '我的分身' },
  { key: '/skill', icon: <ExperimentOutlined />, label: '思维克隆' },
  { key: '/voice', icon: <SoundOutlined />, label: '语音克隆' },
  { key: '/avatar', icon: <PictureOutlined />, label: '虚拟形象' },
];

function HomePage() {
  return (
    <div>
      <Title level={2}>🧬 UploadMyself</Title>
      <Title level={4} type="secondary">
        上传你的数据，生成数字分身
      </Title>
      <div style={{ marginTop: 32 }}>
        {[
          {
            icon: '🧠',
            title: '思维框架克隆',
            desc: '文本语料 → 思维 Skill',
          },
          {
            icon: '🎤',
            title: '语音克隆',
            desc: '语音样本 → 克隆声音',
          },
          {
            icon: '🖼️',
            title: '2D 虚拟形象',
            desc: '一张照片 → 动态说话视频',
          },
          {
            icon: '🧊',
            title: '3D 虚拟形象',
            desc: '一张照片 → 可交互 3D 模型',
          },
        ].map((item) => (
          <div
            key={item.title}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 16,
              padding: '16px 0',
              borderBottom: '1px solid #f0f0f0',
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

function App() {
  const navigate = useNavigate();
  const location = useLocation();

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
        <Content
          style={{
            margin: 24,
            padding: 24,
            background: '#fff',
            borderRadius: 8,
            minHeight: 360,
          }}
        >
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/skill" element={<SkillClone />} />
            <Route path="/voice" element={<VoiceClone />} />
            <Route path="/avatar" element={<AvatarPage />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
}

export default App;
