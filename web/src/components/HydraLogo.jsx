import React from 'react';
import logoImg from '../assets/hydra-logo.jpg';

export default function HydraLogo({ size = 38, glow = true }) {
  return (
    <div 
      className="hydra-logo-wrapper" 
      style={{ 
        display: 'inline-flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        position: 'relative',
        filter: glow ? 'drop-shadow(0 0 8px rgba(255, 0, 60, 0.75)) drop-shadow(0 0 2px rgba(255, 0, 60, 0.9))' : 'none',
        flexShrink: 0
      }}
    >
      <img 
        src={logoImg} 
        alt="HydraForge Cyber Emblem" 
        style={{
          width: `${size}px`,
          height: `${size}px`,
          objectFit: 'cover',
          borderRadius: '4px',
          border: '1.5px solid rgba(255, 0, 60, 0.6)',
          clipPath: 'polygon(0 0, calc(100% - 6px) 0, 100% 6px, 100% 100%, 6px 100%, 0 calc(100% - 6px))'
        }}
      />
    </div>
  );
}
