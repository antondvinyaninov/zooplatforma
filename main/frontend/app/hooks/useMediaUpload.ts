import { useState } from 'react';

export interface UploadedMedia {
  id: number;
  url: string;
  file_name: string;
  original_name: string;
  file_size: number;
  mime_type: string;
  media_type: string;
  width?: number;
  height?: number;
}

export function useMediaUpload() {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const uploadFile = async (file: File, mediaType: string = 'photo'): Promise<UploadedMedia | null> => {
    setUploading(true);
    setProgress(0);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('media_type', mediaType);

      const response = await fetch('http://localhost:8000/api/media/upload', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Upload failed');
      }

      const result = await response.json();
      setProgress(100);
      
      return result.data;
    } catch (error) {
      console.error('Upload error:', error);
      alert(error instanceof Error ? error.message : 'Ошибка загрузки файла');
      return null;
    } finally {
      setUploading(false);
      setTimeout(() => setProgress(0), 1000);
    }
  };

  const uploadMultiple = async (files: File[], mediaType: string = 'photo'): Promise<UploadedMedia[]> => {
    const results: UploadedMedia[] = [];
    
    for (const file of files) {
      const result = await uploadFile(file, mediaType);
      if (result) {
        results.push(result);
      }
    }
    
    return results;
  };

  return {
    uploadFile,
    uploadMultiple,
    uploading,
    progress,
  };
}
