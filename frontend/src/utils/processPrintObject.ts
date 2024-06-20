export const processPrintObject = (data: any, fromPayment: boolean = false) => {
  const tovarCheckObj = data.tovarCheck.map((tovar: any) => {
    const tovarObj = {
      ...tovar,
      name: tovar.tovar_name,
    };
    delete tovarObj.tovar_name;
    return tovarObj;
  });
  const printObj = {
    ...data,
    tovarCheck: tovarCheckObj,
    worker: data.worker_id,
    ...(fromPayment ? { tisCheckUrl: data.link, workerName: "" } : {}),
  };
  delete printObj.worker_id;
  printObj.sum-=printObj.discount
  printObj.discount_percent*=100
  return printObj;
};
