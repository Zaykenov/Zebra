import React, { useEffect, useState } from "react";
import { NextPage } from "next";

import InstructionCarousel, {
  InstructionSlide,
} from "@common/InstructionCarousel";
import { ArrowLeftOnRectangleIcon } from "@heroicons/react/24/solid";
import { useRouter } from "next/router";

const slides: InstructionSlide[] = [
  {
    introText: "Инструкция 1: Начало",
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Для того, чтобы начать работу с ZEBRA CRM необходимо добавить свое
        <br />
        заведение в базу данных. Сделать это можно в навигационной панели
        <br />
        слева в разделе Доступ.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Если ваше заведение уже есть в списке, нажмите кнопку сверху слева для
        завершения
      </p>
    ),
    slideImage: {
      slideImagePath: "1_1.png",
      width: 1100,
      height: 500,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Нажмите на кнопку Добавить и начинайте заполнять данные своего
        <br /> заведения. Учтите, что вам необходимо будет получить Wipon токен
        для <br />
        интеграции с Prosklad.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">Узнать, как получить Wipon токен</p>
    ),
    slideImage: {
      slideImagePath: "1_2.png",
      width: 1100,
      height: 500,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Здесь же можно создать новый счет, склад или сотрудников. Если они
        <br /> уже есть, можно привязать их заведению, которое вы создаете.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Логин и пароль сотрудников позже можно будет просмотреть во вкладке
        Сотрудники раздела Доступ
      </p>
    ),
    slideImage: {
      slideImagePath: "1_3.png",
      width: 1100,
      height: 500,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Вот и все, можно сохранять данные. Вы готовы к дальнейшей работе с
        <br /> ZEBRA CRM.
      </p>
    ),
    slideImage: {
      slideImagePath: "1_4.png",
      width: 1100,
      height: 300,
    },
  },
  {
    introText: "Инструкция 2: Меню",
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        В разделе Меню уже есть позиции, внедренные в меню франчайзером. <br />{" "}
        Вы можете добавить свои собственные, менять цены и состав. Эти <br />
        изменения отразятся только на вашем меню.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Однако стоит учесть, что франчайзер может создавать новые позиции,
        которые окажутся в меню всех франчайзи. <br /> Их тоже можно
        редактировать или удалить
      </p>
    ),
    slideImage: {
      slideImagePath: "2_1.png",
      width: 1200,
      height: 400,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Внедренные франчайзером ингредиенты, тех. карты, товары и их <br />{" "}
        категории можно просмотреть в соответсвующих вкладках. В них же <br />{" "}
        можно добавлять новые категории, ингредиенты, тех. карты и товары.{" "}
        <br /> Они будут только в вашем меню.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Однако стоит учесть, что франчайзер может создавать новые позиции,
        которые окажутся в меню всех франчайзи. <br /> Их тоже можно
        редактировать или удалить
      </p>
    ),
    slideImage: {
      slideImagePath: "2_2.png",
      width: 1100,
      height: 430,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        У тех. карт есть модификаторы. Модификаторы -- это дополнительные <br />{" "}
        ингредиенты, которые гость по своему усмотрению может добавить в <br />{" "}
        свой напиток (сиропы, эссенции, сахар). Просмотреть уже добавленные{" "}
        <br /> или доступные модификаторы тех. карты можно нажав на кнопку
        Редактировать.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Здесь же можно добавить новый модификатор в уже существующий набор или создать совершенно новый <br/> набор модификаторов 
      </p>
    ),
    slideImage: {
      slideImagePath: "2_3.png",
      width: 1000,
      height: 400,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        У тех. карт есть модификаторы. Модификаторы -- это дополнительные <br />{" "}
        ингредиенты, которые гость по своему усмотрению может добавить в <br />{" "}
        свой напиток (сиропы, эссенции, сахар). Просмотреть уже добавленные{" "}
        <br /> или доступные модификаторы тех. карты можно нажав на кнопку
        Редактировать.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Здесь же можно добавить новый модификатор в уже существующий набор или создать совершенно новый <br/> набор модификаторов 
      </p>
    ),
    slideImage: {
      slideImagePath: "2_4.png",
      width: 1000,
      height: 400,
    },
  },
  {
    introText: "Инструкция 3: Склад",
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Добавьте поставки товаров и ингредиентов, чтобы начать продажи. Это <br/> можно сделать как в разделе Склад во вкладке Поставки, так и с <br/> терминала для кассиров. 
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        Поставки можно удалить или отредактировать. Поставки добавленные с терминала тоже.
      </p>
    ),
    slideImage: {
      slideImagePath: "3_1.png",
      width: 1100,
      height: 400,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Количество поставленных товаров начнет отображаться в остатках. По <br/> мере процесса продаж, остатки будут автоматически обновляться и вы <br/> сможете наблюдать за количеством продуктов в складах.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        
      </p>
    ),
    slideImage: {
      slideImagePath: "3_2.png",
      width:1100,
      height: 450,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Количество поставленных товаров начнет отображаться в остатках. По <br/> мере процесса продаж, остатки будут автоматически обновляться и вы <br/> сможете наблюдать за количеством продуктов в складах.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        
      </p>
    ),
    slideImage: {
      slideImagePath: "3_3.png",
      width:1300,
      height: 550,
    },
  },
  {
    headerComponentElement: (
      <p className="text-xl font-semibold">
        Кассиры могут отправить запрос на списание товара или ингредиента с <br/> терминала, которые после можно отклонить или принять в админ. <br/> панели.
      </p>
    ),
    footerComponentElement: (
      <p className="text-xl font-semibold">
        
      </p>
    ),
    slideImage: {
      slideImagePath: "3_4.png",
      width:1400,
      height: 550,
    },
  },
  {
    introText: "Поздравляем, вы прошли инструкции",
  },
];

const GuidePage: NextPage = () => {
  const router = useRouter();
  return (
    <div className="relative">
      <button className="w-20 h-20" onClick={() => router.push("/menu")}>
        <ArrowLeftOnRectangleIcon
          className="text-gray-400 absolute top-0 left-0 mt-4 ml-4"
          width={40}
          height={40}
        />
      </button>
      <InstructionCarousel slides={slides} />
    </div>
  );
};

export default GuidePage;
