import { useState } from "react";

const useFetching = (callback: (...args: any[]) => Promise<any>) => {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");

    const fetchData = async (...args: any[]) => {
        try {
            setIsLoading(true);
            await callback(...args);
        } catch (e: any) {
            setError(e.message);
        } finally {
            setIsLoading(false);
        }
    };

    return {fetchData, isLoading, error}; 
};

export default useFetching