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
    console.log('🔄 uploadFile начат для:', file.name);
    setUploading(true);
    setProgress(0);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('media_type', mediaType);
      
      console.log('📤 Отправка запроса на http://localhost:8000/api/media/upload');

      const response = await fetch('http://localhost:8000/api/media/upload', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      console.log('📥 Ответ получен, статус:', response.status);

      if (!response.ok) {
        const error = await response.json();
        console.error('❌ Ошибка от сервера:', error);
        throw new Error(error.error || 'Upload failed');
      }

      const result = await response.json();
      console.log('✅ Результат загрузки:', result);
      setProgress(100);
      
      return result.data;
    } catch (error) {
      console.error('❌ Upload error:', error);
      alert(error instanceof Error ? error.message : 'Ошибка загрузки файла');
      return null;
    } finally {
      setUploading(false);
      setTimeout(() => setProgress(0), 1000);
    }
  };

  const uploadMultiple = async (files: File[], mediaType: string = 'photo'): Promise<UploadedMedia[]> => {
    console.log('📦 uploadMultiple начат для', files.length, 'файлов');
    const results: UploadedMedia[] = [];
    
    for (const file of files) {
      console.log('⏳ Загрузка файла:', file.name);
      const result = await uploadFile(file, mediaType);
      if (result) {
        console.log('✅ Файл загружен:', result);
        results.push(result);
      } else {
        console.log('❌ Файл не загружен:', file.name);
      }
    }
    
    console.log('🎉 uploadMultiple завершен, загружено:', results.length);
    return results;
  };

  return {
    uploadFile,
    uploadMultiple,
    uploading,
    progress,
  };
}
