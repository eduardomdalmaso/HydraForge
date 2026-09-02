import React from 'react';

export default function DatasetCyberIcon({ name, color = 'currentColor', size = 16 }) {
  const props = { width: size, height: size, viewBox: '0 0 24 24', fill: 'none', stroke: color, strokeWidth: 2, strokeLinecap: 'square', strokeLinejoin: 'miter' };

  switch (name) {
    case 'zip':
      return (
        <svg {...props}>
          <path d="M4 4h16v16H4z" />
          <path d="M10 4v4h4V4M10 8h4v4h-4M10 12v4h4v-4M12 16v2" />
        </svg>
      );
    case 'tag':
      return (
        <svg {...props}>
          <path d="M12 2l10 10-9 9L3 11V2h9z" />
          <circle cx="7.5" cy="6.5" r="1.5" fill={color} />
        </svg>
      );
    case 'fusion':
      return (
        <svg {...props}>
          <path d="M4 6h6l4 6h6M4 18h6l4-6" />
          <circle cx="20" cy="12" r="2" fill={color} />
          <circle cx="4" cy="6" r="2" fill={color} />
          <circle cx="4" cy="18" r="2" fill={color} />
        </svg>
      );
    case 'kfold':
      return (
        <svg {...props}>
          <path d="M12 2v20M2 12h20M4.93 4.93l14.14 14.14M4.93 19.07L19.07 4.93" strokeDasharray="3 3" />
          <circle cx="12" cy="12" r="9" />
        </svg>
      );
    case 'shield':
      return (
        <svg {...props}>
          <path d="M12 2L3 6v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V6l-9-4z" />
          <path d="M9 12l2 2 4-4" />
        </svg>
      );
    case 'launch':
      return (
        <svg {...props}>
          <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" fill={color} fillOpacity="0.2" />
        </svg>
      );
    case 'copy':
      return (
        <svg {...props}>
          <rect x="8" y="8" width="13" height="13" rx="1" />
          <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
        </svg>
      );
    case 'car':
      return (
        <svg {...props}>
          <path d="M5 12h14l-1.5-5H6.5L5 12zM3 12h18v5H3v-5z" />
          <circle cx="7" cy="17" r="2" fill={color} />
          <circle cx="17" cy="17" r="2" fill={color} />
        </svg>
      );
    case 'person':
      return (
        <svg {...props}>
          <circle cx="12" cy="7" r="4" />
          <path d="M5 21v-2a7 7 0 0 1 14 0v2" />
        </svg>
      );
    default:
      return (
        <svg {...props}>
          <polygon points="12 2 22 12 12 22 2 12" />
        </svg>
      );
  }
}
