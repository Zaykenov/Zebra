namespace ExcelGenApi.Models
{
    public record NeedToCreateExcelReq(string fileName, string[] headersRow, string[][] dataRows);
}