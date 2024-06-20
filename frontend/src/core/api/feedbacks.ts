import { QueryOptionsData, backendApiInstance, queryOptionsToString } from ".";

export const getAllFeedbacks = async (queryOptions?: QueryOptionsData) => {
  try {
    const response = await backendApiInstance(
      "/mobile/feedback/getAll" + queryOptionsToString(queryOptions, false)
    );
    const users = response.data;
    return users;
  } catch (error) {
    console.error("Error fetching feedbacks:", error);
    throw error;
  }
};

export const getFeedbackById = async (id: number) => {
  try {
    const response = await backendApiInstance(`/mobile/feedback/get/${id}`);
    const feedback = response.data;
    return feedback;
  } catch (error) {
    console.error("Error fetching feedback by id:", error);
    throw error;
  }
};
