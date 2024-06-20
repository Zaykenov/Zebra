using ExcelGenApi.Interfaces;
using ExcelGenApi.Models;
using Microsoft.AspNetCore.Mvc;

namespace ExcelGenApi.Controllers {
    [ApiController]
    [Route("excel-gen")]
    public class ExcelGenController : ControllerBase {

        private readonly IExcelGenerator _service;

        public ExcelGenController(IExcelGenerator service) {
            _service = service;
        }

        [HttpPost]
        [Route("generate")]
        public async Task<FileResult> GenerateBaseExcel(NeedToCreateExcelReq req) {
            var genRes =  await _service.GenerateBaseExcel(req);
            return File(genRes.Content, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileDownloadName: genRes.FileName);
        }
        [HttpPost]
        [Route("generateInvent")]
        public async Task<FileResult> GeneratePartialInvent(PartialInventExcelReq req) {
            var genRes =  await _service.GeneratePartialInventExcel(req);
            return File(genRes.Content, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileDownloadName: genRes.FileName);
        }
    }
}
