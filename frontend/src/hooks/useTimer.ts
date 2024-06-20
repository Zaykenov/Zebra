import { useEffect, useState } from "react";

const useTimer = (initialSeconds: number) => {
  const [seconds, setSeconds] = useState(initialSeconds);

  const [isExpired, setIsExpired] = useState<boolean>(false);
  useEffect(() => {
    const interval = setInterval(() => {
      setSeconds((prevSeconds) => {
        const newSeconds = prevSeconds > 0 ? prevSeconds - 1 : 0;
        setIsExpired(newSeconds === 0);
        return newSeconds;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return { seconds, isExpired };
};

export default useTimer;
