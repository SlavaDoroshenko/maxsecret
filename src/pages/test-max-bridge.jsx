import { useEffect, useState } from "react";

function Test() {
  const [webApp, setWebApp] = useState(null);
  const [userData, setUserData] = useState(null);
  const [platformData, setPlatformData] = useState(null);
  const [hapticFeedback, setHapticFeedback] = useState(null);

  useEffect(() => {
    if (window.WebApp) {
      window.WebApp.ready();
      const app = window.WebApp;
      setWebApp(app);
      setUserData(app.initData);
      setPlatformData(app.platform); // ✅ теперь app — это window.WebApp
      setHapticFeedback(app.HapticFeedback);
    }
  }, []);

  if (!webApp) {
    return <div>Загрузка...</div>;
  }

  return (
    <div>
      <h1>Привет, {userData}!</h1>
      <h1>{platformData}</h1>
      <button
        onClick={() => {
          webApp.HapticFeedback.impactOccurred("heavy", false);
          webApp.openCodeReader();
        }}
      >
        Закрыть
      </button>
    </div>
  );
}

export default Test;
