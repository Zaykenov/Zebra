import http from 'k6/http';

export const options = {
  scenarios: {
    my_scenario1: {
      executor: 'constant-arrival-rate',
      duration: '30s', // total duration
      preAllocatedVUs: 50, // to allocate runtime resources     preAll

      rate: 50, // number of constant iterations given `timeUnit`
      timeUnit: '1s',
    },
  },
};

export default function () {
  const payload = JSON.stringify({
    "discount_percent": 0,
    "payment": "Картой",
    "cash": 0,
    "card": 2280,
    "comment": "Пейджер №1: ",
    "status": "closed",
    "mobile_user_id": null,
    "tovarCheck": [
      {
        "tovar_id": 27,
        "quantity": 1,
        "comments": "",
        "modifications": ""
      }
    ],
    "techCartCheck": [
      {
        "tech_cart_id": 103,
        "quantity": 1,
        "comments": "",
        "modificators": [
          {
            "id": 140,
            "nabor_id": 141,
            "name": "Миндальное молоко",
            "quantity": 1,
            "totalBrutto": 0.3,
            "brutto": 0.3,
            "cost": 0,
            "price": 350
          }
        ]
      }
    ]
  });
  const headers = { 
    'Content-Type': 'application/json',
    'Authorization':'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTE4OTM3NTUsImlhdCI6MTY5MTgxODE1NSwicm9sZSI6IndvcmtlciIsInVzZXJfaWQiOjE1LCJzaG9wcyI6WzRdLCJiaW5kX3Nob3AiOjR9.5ItF3wXuNUM4N1ok8IuAYTFZAoSfRdsKTwAutllW9PE',
    'Idempotency-Key': '9044bbc6-1b24-3416-a9d4-025a75e887ed_3'
   };
  http.post('https://zebra-api.korsetu.kz/check/create', payload, { headers });
}
