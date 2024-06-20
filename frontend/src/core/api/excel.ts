import { excelGeneratorApiInstance } from ".";
import downloadExcelFile from "@utils/downloadExcelFile";
const generateDataRows = (postArray: any[]) => {
    return postArray.map((elem) => {
      return Object.values(elem);
    });
  };

export const getExcelFile = async (fileName: string, headersRow: any[], excelData: any[]) => {
    const dataRows = generateDataRows(excelData)
    const response = await excelGeneratorApiInstance.post(
      "/excel-gen/generate",
      {
        fileName,
        headersRow,
        dataRows
      },
      {
        responseType: 'arraybuffer', 
      }
    );
    const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    downloadExcelFile(blob, fileName)
};

export const getParitalInventExcelFile = async (fileName: string, excelData: any[][]) => {
  const response = await excelGeneratorApiInstance.post(
    "/excel-gen/generateInvent",
    {
      fileName,
      excelData
    },
    {
      responseType: 'arraybuffer', 
    }
  );
  const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  downloadExcelFile(blob, fileName)
};