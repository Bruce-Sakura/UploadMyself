import { useRef, useEffect, useState } from 'react';
import { Button, Typography } from 'antd';
import { CaretRightOutlined, PauseOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface Joint {
  x: number;
  y: number;
}

interface AnimationData {
  idle: { frame: number; joints: Record<string, Joint> }[];
  wave: { frame: number; joints: Record<string, Joint> }[];
  fps: number;
  skeleton: {
    joints: Record<string, Joint>;
    bones: [string, string][];
  };
}

interface Props {
  animationData: AnimationData;
  cartoonImageUrl?: string;
  width?: number;
  height?: number;
}

export default function SkeletonCanvas({ animationData, cartoonImageUrl, width = 400, height = 500 }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [playing, setPlaying] = useState(false);
  const [animType, setAnimType] = useState<'idle' | 'wave'>('idle');
  const frameRef = useRef(0);
  const timerRef = useRef<number | null>(null);
  const bgImageRef = useRef<HTMLImageElement | null>(null);

  // Load background image
  useEffect(() => {
    if (cartoonImageUrl) {
      const img = new Image();
      img.crossOrigin = 'anonymous';
      img.onload = () => {
        bgImageRef.current = img;
        drawFrame(0, animType);
      };
      img.src = cartoonImageUrl;
    }
  }, [cartoonImageUrl]);

  const drawFrame = (frameIdx: number, type: 'idle' | 'wave') => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const anim = animationData[type];
    if (!anim || anim.length === 0) return;

    const frame = anim[frameIdx % anim.length];
    const joints = frame.joints;
    const bones = animationData.skeleton.bones;

    // Clear canvas
    ctx.clearRect(0, 0, width, height);

    // Draw background image (scaled to fit)
    if (bgImageRef.current) {
      const img = bgImageRef.current;
      const scale = Math.min(width / img.width, height / img.height);
      const w = img.width * scale;
      const h = img.height * scale;
      ctx.drawImage(img, (width - w) / 2, (height - h) / 2, w, h);
    }

    // Scale joints to canvas size
    const skelJoints = animationData.skeleton.joints;
    const maxX = Math.max(...Object.values(skelJoints).map(j => j.x), 1);
    const maxY = Math.max(...Object.values(skelJoints).map(j => j.y), 1);
    const scaleX = width / (maxX * 1.5);
    const scaleY = height / (maxY * 1.2);

    const scaleJoint = (j: Joint): Joint => ({
      x: j.x * scaleX,
      y: j.y * scaleY,
    });

    // Draw bones
    ctx.strokeStyle = '#00ff88';
    ctx.lineWidth = 3;
    ctx.lineCap = 'round';
    for (const [j1Name, j2Name] of bones) {
      const j1 = joints[j1Name];
      const j2 = joints[j2Name];
      if (j1 && j2) {
        const s1 = scaleJoint(j1);
        const s2 = scaleJoint(j2);
        ctx.beginPath();
        ctx.moveTo(s1.x, s1.y);
        ctx.lineTo(s2.x, s2.y);
        ctx.stroke();
      }
    }

    // Draw joints
    for (const [name, j] of Object.entries(joints)) {
      const s = scaleJoint(j);
      ctx.beginPath();
      ctx.arc(s.x, s.y, name === 'head' ? 8 : 4, 0, Math.PI * 2);
      ctx.fillStyle = name === 'head' ? '#ff4444' : '#4488ff';
      ctx.fill();
      ctx.strokeStyle = '#fff';
      ctx.lineWidth = 1;
      ctx.stroke();

      // Label
      ctx.fillStyle = '#fff';
      ctx.font = '10px sans-serif';
      ctx.fillText(name.slice(0, 3), s.x + 8, s.y - 4);
    }

    // Frame info
    ctx.fillStyle = 'rgba(0,0,0,0.5)';
    ctx.fillRect(0, height - 24, width, 24);
    ctx.fillStyle = '#fff';
    ctx.font = '12px sans-serif';
    ctx.fillText(`${type} | Frame ${frameIdx + 1}/${anim.length}`, 8, height - 8);
  };

  const startAnimation = () => {
    if (timerRef.current) return;
    setPlaying(true);
    const anim = animationData[animType];
    const fps = animationData.fps || 30;

    const tick = () => {
      frameRef.current = (frameRef.current + 1) % anim.length;
      drawFrame(frameRef.current, animType);
      timerRef.current = window.setTimeout(tick, 1000 / fps);
    };
    timerRef.current = window.setTimeout(tick, 1000 / fps);
  };

  const stopAnimation = () => {
    setPlaying(false);
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  };

  const switchAnim = (type: 'idle' | 'wave') => {
    stopAnimation();
    setAnimType(type);
    frameRef.current = 0;
    drawFrame(0, type);
  };

  // Draw initial frame
  useEffect(() => {
    drawFrame(0, animType);
  }, [animationData]);

  return (
    <div>
      <canvas
        ref={canvasRef}
        width={width}
        height={height}
        style={{
          border: '1px solid #f0f0f0',
          borderRadius: 8,
          background: '#1a1a2e',
          display: 'block',
        }}
      />
      <div style={{ marginTop: 8, display: 'flex', gap: 8, alignItems: 'center' }}>
        <Button
          type="primary"
          icon={playing ? <PauseOutlined /> : <CaretRightOutlined />}
          onClick={playing ? stopAnimation : startAnimation}
          size="small"
        >
          {playing ? '暂停' : '播放'}
        </Button>
        <Button
          size="small"
          type={animType === 'idle' ? 'default' : 'dashed'}
          onClick={() => switchAnim('idle')}
        >
          待机
        </Button>
        <Button
          size="small"
          type={animType === 'wave' ? 'default' : 'dashed'}
          onClick={() => switchAnim('wave')}
        >
          挥手
        </Button>
        <Text type="secondary" style={{ fontSize: 11, marginLeft: 'auto' }}>
          15个关节 · 14根骨骼
        </Text>
      </div>
    </div>
  );
}
