import { Metadata } from 'next';
import HomeClient from './HomeClient';

type Props = {
  searchParams: Promise<{ metka?: string }>;
};

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const params = await searchParams;
  const metkaId = params.metka;

  if (metkaId) {
    try {
      // Получаем данные поста для SEO
      const response = await fetch(`http://localhost:8000/api/posts/${metkaId}`, {
        cache: 'no-store',
      });

      if (response.ok) {
        const result = await response.json();
        const post = result.data;

        // Формируем описание из контента поста
        const description = post.content.substring(0, 160) + (post.content.length > 160 ? '...' : '');
        
        // Получаем имя автора
        const authorName = post.author_type === 'organization' 
          ? (post.organization?.short_name || post.organization?.name || 'Организация')
          : (post.user?.name || 'Пользователь') + (post.user?.last_name ? ' ' + post.user.last_name : '');

        // Получаем первое изображение если есть
        const image = post.attachments?.find((a: any) => a.type === 'image' || a.media_type === 'image')?.url;

        return {
          title: `${authorName}: ${description}`,
          description: description,
          openGraph: {
            title: `${authorName}: ${description}`,
            description: description,
            images: image ? [`http://localhost:8000${image}`] : [],
            type: 'article',
            url: `http://localhost:3000/?metka=${metkaId}`,
          },
          twitter: {
            card: 'summary_large_image',
            title: `${authorName}: ${description}`,
            description: description,
            images: image ? [`http://localhost:8000${image}`] : [],
          },
        };
      }
    } catch (error) {
      console.error('Error fetching post for SEO:', error);
    }
  }

  // Дефолтные meta-теги для главной страницы
  return {
    title: 'Главная - Зооплатформа',
    description: 'Социальная сеть для владельцев домашних животных',
  };
}

export default async function Home({ searchParams }: Props) {
  const params = await searchParams;
  return <HomeClient searchParams={params} />;
}
