import React, {
  Context,
  createContext,
  FC,
  ReactNode,
  useContext,
  useEffect,
} from "react";
import { useRouter } from "next/router";
import { WorkerRole } from "../core/api/workers";
interface AuthContextProps {
  authState?: {
    token: string;
    access: 1 | 2 | null;
  };
  setAuthInfo?: (data: any) => void;
  isAuthenticated?: () => boolean;
  logOut?: () => void;
}

export const AuthContext: Context<AuthContextProps> = createContext({});
const { Provider } = AuthContext;

const AuthProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const router = useRouter();

  useEffect(() => {
    if (isAuthenticated()) {
      const role = getRole();
      if (router.pathname === "/login") {
        role === WorkerRole.WORKER && router.push("/terminal/order");
        (role === WorkerRole.MANAGER || role === WorkerRole.MASTER) &&
          router.push("/menu");
      } else {
        if (
          role === WorkerRole.WORKER &&
          !router.pathname.startsWith("/terminal")
        ) {
          router.push("/terminal/order");
        }
        if (
          role === WorkerRole.MANAGER &&
          router.pathname.startsWith("/terminal")
        ) {
          router.push("/menu");
        }
      }
    } else {
      if (router.pathname !== "/login" && router.pathname === "privacy") {
        router.push("/login");
      }
    }
  }, [router]);

  const setAuthInfo = (data: any) => {
    localStorage.setItem("zebra.authToken", data.token);
    localStorage.setItem("zebra.role", data.role);
  };

  const isAuthenticated = () => {
    return !!localStorage.getItem("zebra.authToken");
  };

  const getRole = () => {
    return localStorage.getItem("zebra.role");
  };

  const logOut = () => {
    localStorage.removeItem("zebra.authToken");
    localStorage.removeItem("zebra.role");
  };

  return (
    <Provider value={{ setAuthInfo, isAuthenticated, logOut }}>
      {children}
    </Provider>
  );
};

export const useAuth = (): AuthContextProps => {
  return useContext(AuthContext);
};

export default AuthProvider;
