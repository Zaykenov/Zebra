using ExcelGenApi.Models;

namespace ExcelGenApi.Interfaces
{
    public interface IExcelGenerator {
        public Task<GenerateBaseExcelRes> GenerateBaseExcel(NeedToCreateExcelReq req);

        public Task<GenerateBaseExcelRes> GeneratePartialInventExcel(PartialInventExcelReq req);
    }
} 