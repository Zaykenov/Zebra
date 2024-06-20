import type { NextPage } from "next";
import PageLayout from "../components/__layouts/PageLayout";
import { useEffect } from "react";
import { useRouter } from "next/router";

const Home: NextPage = () => {
  const router = useRouter();
  useEffect(() => {
    router.replace("/login");
  }, [router]);

  return (
    <div className="">
      <PageLayout />
    </div>
  );
};

export default Home;
