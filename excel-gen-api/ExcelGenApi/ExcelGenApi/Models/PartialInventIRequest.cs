namespace ExcelGenApi.Models
{
    public record PartialInventItem{string name; int quantityDifference; int sumDifference;}
    public record PartialInventExcelReq(string fileName, object[][][] excelData);
}