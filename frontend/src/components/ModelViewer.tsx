import { Suspense, useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, useGLTF, Environment, ContactShadows } from '@react-three/drei';
import { Spin, Typography, Button } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';

const { Text } = Typography;

function Model({ url }: { url: string }) {
  const { scene } = useGLTF(url);
  const ref = useRef<any>();

  useFrame((state) => {
    if (ref.current) {
      ref.current.rotation.y = Math.sin(state.clock.elapsedTime * 0.5) * 0.1;
    }
  });

  return <primitive ref={ref} object={scene} scale={1} />;
}

function LoadingSpinner() {
  return (
    <div style={{
      position: 'absolute',
      top: '50%',
      left: '50%',
      transform: 'translate(-50%, -50%)',
      textAlign: 'center',
    }}>
      <Spin size="large" />
      <div style={{ marginTop: 16, color: '#888' }}>加载 3D 模型中...</div>
    </div>
  );
}

interface Props {
  modelUrl: string;
  width?: number;
  height?: number;
  autoRotate?: boolean;
  showDownload?: boolean;
}

export default function ModelViewer({
  modelUrl,
  width = 400,
  height = 500,
  autoRotate = true,
  showDownload = true,
}: Props) {
  return (
    <div>
      <div style={{
        width,
        height,
        borderRadius: 8,
        overflow: 'hidden',
        background: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)',
        position: 'relative',
      }}>
        <Suspense fallback={<LoadingSpinner />}>
          <Canvas
            camera={{ position: [0, 1, 3], fov: 45 }}
            style={{ width, height }}
          >
            <ambientLight intensity={0.5} />
            <directionalLight position={[5, 5, 5]} intensity={1} />
            <directionalLight position={[-5, 5, -5]} intensity={0.5} />

            <Model url={modelUrl} />

            <ContactShadows
              position={[0, -0.5, 0]}
              opacity={0.4}
              scale={5}
              blur={2}
            />

            <Environment preset="city" />

            <OrbitControls
              autoRotate={autoRotate}
              autoRotateSpeed={2}
              enablePan={true}
              enableZoom={true}
              minDistance={1}
              maxDistance={10}
              target={[0, 0.5, 0]}
            />
          </Canvas>
        </Suspense>

        {/* 控件提示 */}
        <div style={{
          position: 'absolute',
          bottom: 8,
          left: 8,
          background: 'rgba(0,0,0,0.5)',
          borderRadius: 4,
          padding: '2px 8px',
        }}>
          <Text style={{ color: '#fff', fontSize: 11 }}>
            🖱️ 拖拽旋转 · 滚轮缩放
          </Text>
        </div>
      </div>

      {/* 操作按钮 */}
      <div style={{ marginTop: 8, display: 'flex', gap: 8 }}>
        {showDownload && (
          <Button
            icon={<DownloadOutlined />}
            size="small"
            href={modelUrl}
            download
          >
            下载 GLB
          </Button>
        )}
        <Text type="secondary" style={{ fontSize: 11, alignSelf: 'center' }}>
          Three.js 实时渲染 · WebGL
        </Text>
      </div>
    </div>
  );
}
