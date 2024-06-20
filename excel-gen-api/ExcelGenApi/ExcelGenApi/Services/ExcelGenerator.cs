using System.Drawing;
using ExcelGenApi.Interfaces;
using ExcelGenApi.Models;
using OfficeOpenXml;
using OfficeOpenXml.Style;

namespace ExcelGenApi.Services {
    public class ExcelGenerator : IExcelGenerator {
        public ExcelGenerator() {
        }


        private static object TryConvertToNumber(string value)
        {
            if (double.TryParse(value, out double result))
            {
                return result;
            }
            return value;
        }

        private static int[] GetLongestStringLengths(NeedToCreateExcelReq req) {
            var combinedRows = new[] { req.headersRow }.Concat(req.dataRows).ToArray();
            var longestStringLengths = new int[req.headersRow.Length];

            for (int i = 0; i < combinedRows.Length; i++) {
                for (int j = 0; j < combinedRows[i].Length; j++) {
                    var valueLength = combinedRows[i][j].Length;
                    if (valueLength > longestStringLengths[j]) {
                        longestStringLengths[j] = valueLength;
                    }
                }
            }

            return longestStringLengths;
        }

        public async Task<GenerateBaseExcelRes> GenerateBaseExcel(NeedToCreateExcelReq req) {
            int[] longestStringSizes = GetLongestStringLengths(req);

            string fileName = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString() + ".xlsx");
            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;

            using (ExcelPackage package = new ExcelPackage()) {
                ExcelWorksheet worksheet = package.Workbook.Worksheets.Add("result");

                for (var i = 0; i < req.headersRow.Length; i++) {
                    string headerCell = req.headersRow[i];
                    worksheet.Column(i + 1).Width = longestStringSizes[i];
                    worksheet.Cells[1, i + 1].Value = headerCell;
                }

                var headerRange = worksheet.Cells[1, 1, 1, req.headersRow.Length];
                headerRange.Style.Fill.PatternType = ExcelFillStyle.Solid;
                headerRange.Style.Fill.BackgroundColor.SetColor(Color.LightBlue);

                for (var i = 0; i < req.dataRows.Length; i++) {
                    var row = req.dataRows[i];
                    var cells = row;
                    for (var j = 0; j < cells.Length; j++) {
                        var value = cells[j];
                        worksheet.Cells[2 + i, j + 1].Value = TryConvertToNumber(value);
                    }
                }

                await package.SaveAsAsync(new FileInfo(fileName));

                return new GenerateBaseExcelRes(await File.ReadAllBytesAsync(fileName), req.fileName);
            }
        }


        public async Task<GenerateBaseExcelRes> GeneratePartialInventExcel(PartialInventExcelReq req) {
            string fileName = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString() + ".xlsx");
            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;
            using (var package = new ExcelPackage())
            {
                ExcelWorksheet worksheet = package.Workbook.Worksheets.Add("Таблицы");

                int currentRow = 1;
                int currentColumn = 1;

                foreach (var table in req.excelData)
                {
                    foreach (var row in table)
                    {
                        for (int i = 0; i < row.Count(); i++)
                        {
                            ExcelRange cell = worksheet.Cells[currentRow, currentColumn + i];
                            cell.Value = TryConvertToNumber(row[i].ToString());

                            if (currentRow == 1)
                            {
                                var style = cell.Style;
                                style.Font.Bold = true;
                                style.HorizontalAlignment = ExcelHorizontalAlignment.Center;
                                style.VerticalAlignment = ExcelVerticalAlignment.Center;
                                style.Fill.PatternType = ExcelFillStyle.Solid;
                                style.Fill.BackgroundColor.SetColor(System.Drawing.Color.LightGray);
                            }
                        }

                        currentRow++;
                    }

                    int startColumn = currentColumn;
                    int endColumn = currentColumn + table[0].Count() - 1;
                    worksheet.Cells[1, startColumn, currentRow - 1, endColumn].AutoFitColumns();

                    currentColumn += table[0].Count() + 1;
                    currentRow = 1;
                }

                await package.SaveAsAsync( new FileInfo(fileName));
                return new GenerateBaseExcelRes(await File.ReadAllBytesAsync(fileName), req.fileName);

            }
        }
    }
}
