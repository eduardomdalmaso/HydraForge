import { useEffect, useRef } from 'react';

export function useWebcamStream(isWebcam) {
  const videoRef = useRef(null);
  const streamTrackRef = useRef(null);

  useEffect(() => {
    if (isWebcam && navigator.mediaDevices?.getUserMedia) {
      navigator.mediaDevices.getUserMedia({ video: { width: 640, height: 480 } })
        .then(s => {
          streamTrackRef.current = s;
          if (videoRef.current) videoRef.current.srcObject = s;
        })
        .catch(err => console.warn('Webcam permission:', err));
    } else if (streamTrackRef.current) {
      streamTrackRef.current.getTracks().forEach(t => t.stop());
      streamTrackRef.current = null;
    }
    return () => streamTrackRef.current?.getTracks().forEach(t => t.stop());
  }, [isWebcam]);

  const captureFrame = () => {
    if (isWebcam && videoRef.current && videoRef.current.videoWidth) {
      const canvas = document.createElement('canvas');
      canvas.width = videoRef.current.videoWidth;
      canvas.height = videoRef.current.videoHeight;
      canvas.getContext('2d').drawImage(videoRef.current, 0, 0);
      return canvas.toDataURL('image/jpeg', 0.85);
    }
    return null;
  };

  return { videoRef, captureFrame };
}
