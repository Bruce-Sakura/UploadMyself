import { Routes, Route } from "react-router-dom";
import { Layout } from "antd";
import {
  HomeOutlined,
  SoundOutlined,
  PictureOutlined,
  BoxPlotOutlined,
  ExperimentOutlined,
} from "@ant-design/icons";

const { Header, Content, Sider } = Layout;

const menuItems = [
  { key: "/", icon: <HomeOutlined />, label: "首页" },
  { key: "/skill", icon: <ExperimentOutlined />, label: "思维克隆" },
  { key: "/voice", icon: <SoundOutlined />, label: "语音克隆" },
  { key: "/avatar-2d", icon: <PictureOutlined />, label: "2D 形象" },
  { key: "/avatar-3d", icon: <BoxPlotOutlined />, label: "3D 形象" },
];

function App() {
  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider theme="dark" width={200}>
        <div
          style={{
            height: 64,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "#fff",
            fontSize: 18,
            fontWeight: "bold",
          }}
        >
          🧬 UploadMyself
        </div>
      </Sider>
      <Layout>
        <Header
          style={{
            background: "#fff",
            padding: "0 24px",
            display: "flex",
            alignItems: "center",
          }}
        >
          <h2 style={{ margin: 0 }}>克隆你自己</h2>
        </Header>
        <Content style={{ margin: 24, padding: 24, background: "#fff" }}>
          <Routes>
            <Route
              path="/"
              element={
                <div>
                  <h1>🧬 UploadMyself</h1>
                  <p>上传你的数据，生成数字分身</p>
                  <ul>
                    <li>🧠 思维框架克隆 — 文本语料 → 思维 Skill</li>
                    <li>🎤 语音克隆 — 语音样本 → 克隆声音</li>
                    <li>🖼️ 2D 虚拟形象 — 一张照片 → 动态说话视频</li>
                    <li>🧊 3D 虚拟形象 — 一张照片 → 可交互 3D 模型</li>
                    <li>🔬 模型蒸馏 — 大模型压缩，降低成本</li>
                  </ul>
                </div>
              }
            />
            <Route path="/skill" element={<div>思维框架克隆 (开发中)</div>} />
            <Route path="/voice" element={<div>语音克隆 (开发中)</div>} />
            <Route path="/avatar-2d" element={<div>2D 形象 (开发中)</div>} />
            <Route path="/avatar-3d" element={<div>3D 形象 (开发中)</div>} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
}

export default App;
