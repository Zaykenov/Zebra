import { QueryOptionsData, backendApiInstance, queryOptionsToString } from ".";

export const getAbcAnalysis = async (queryOptions?: QueryOptionsData) => {
    const res = await backendApiInstance.get(
        `statistics/ABC` + queryOptionsToString(queryOptions, false)
    );
    return res.data;
};