import { QueryOptionsData, backendApiInstance, queryOptionsToString } from "@api/index";

class ClientsService { 
  getAllClients = async (queryOptions?: QueryOptionsData) => {
    try {
      const response = await backendApiInstance(
        "/mobile/user/getAll" + queryOptionsToString(queryOptions, false)
      );
      const users = response.data;
      return users;
    } catch (error) {
      console.error("Error fetching mobile users:", error);
      throw error;
    }
  }
  getClientById = async (id: string) => {
    try {
      const response = await backendApiInstance(`/mobile/user/get/${id}`);
      const client = response.data;
      return client;
    } catch (error) {
      console.error("Error fetching client by id:", error);
      throw error;
    }
  };
}

export default new ClientsService()