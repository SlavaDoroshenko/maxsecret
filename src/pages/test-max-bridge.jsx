import { useEffect, useState } from "react";

function Test() {
  const [webApp, setWebApp] = useState(null);
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    // Проверяем, доступен ли WebApp
    if (window.WebApp) {
      // Сообщаем MAX, что приложение готово
      window.WebApp.ready();

      // Сохраняем ссылку для удобства
      setWebApp(window.WebApp);

      setUserData(webApp.initDataUnsafe);

      // Данные пользователя (небезопасные, только для UI!)
      //   window.WebApp.openLink(
      //     "https://www.figma.com/design/HXV56zy3MEmVOfG2jakrdS/Untitled?node-id=0-1&p=f&t=DmX7rmNaLFvTYrSd-0"
      //   );
    }
  }, []);

  if (!webApp) {
    return <div>Загрузка...</div>;
  }

  return (
    <div>
      <h1>Привет, {userData}!</h1>
      <button onClick={() => webApp.close()}>Закрыть</button>
    </div>
  );
}

export default Test;
