CREATE TABLE IF NOT EXISTS category_tovar (
  id bigserial,
  name text NOT NULL,
  image text,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS category_ingredient (
  id bigserial,
  name text NOT NULL,
  image text,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS tovar (
  id bigserial,
  name text NOT NULL ,
  category int NOT NULL,
  image text,
  tax text NOT NULL,
  measure text NOT NULL,
  cost float NOT NULL ,
  price float NOT NULL ,
  profit float NOT NULL ,
  margin int NOT NULL ,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id),
  FOREIGN KEY(category) REFERENCES category_tovar(id)
);
CREATE TABLE IF NOT EXISTS tech_cart (
  id bigserial,
  name text NOT NULL ,
  category int NOT NULL,
  image text,
  tax text NOT NULL,
  measure text NOT NULL,
  cost float NOT NULL ,
  price float NOT NULL ,
  profit float NOT NULL ,
  margin int NOT NULL ,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id),
  FOREIGN KEY(category) REFERENCES category_tovar(id)
);

CREATE TABLE IF NOT EXISTS sklad (
  id bigserial,
  name text NOT NULL ,
  address text NOT NULL,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS ingredient (
  id bigserial ,
  name text NOT NULL ,
  category int NOT NULL,
  image text,
  measure text NOT NULL,
  cost float NOT NULL ,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id),
  FOREIGN KEY(category) REFERENCES category_ingredient(id)
);

CREATE TABLE IF NOT EXISTS tech_cart_ingredient (
  tech_cart_id bigserial ,
  ingredient_id bigserial ,
  brutto float DEFAULT 0,
  PRIMARY KEY(tech_cart_id, ingredient_id),
  FOREIGN KEY(ingredient_id) REFERENCES ingredient(id),
  FOREIGN KEY(tech_cart_id) REFERENCES tech_cart(id)
);

CREATE TABLE IF NOT EXISTS sklad_ingredient (
  sklad_id bigserial,
  ingredient_id bigserial,
  quantity float NOT NULL,
  PRIMARY KEY(sklad_id, ingredient_id),
  FOREIGN KEY(sklad_id) REFERENCES sklad(id),
  FOREIGN KEY(ingredient_id) REFERENCES ingredient(id)
);

CREATE TABLE IF NOT EXISTS sklad_tovar (
  sklad_id bigserial,
  tovar_id bigserial,
  quantity float NOT NULL,
  PRIMARY KEY(sklad_id, tovar_id),
  FOREIGN KEY(sklad_id) REFERENCES sklad(id),
  FOREIGN KEY(tovar_id) REFERENCES tovar(id)
);

CREATE TABLE IF NOT EXISTS nabor (
  id bigserial,
  name text NOT NULL,
  min int,
  max int,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS nabor_ingredient (
  nabor_id bigserial,
  ingredient_id bigserial,
  brutto float NOT NULL,
  price float NOT NULL,
  PRIMARY KEY(nabor_id, ingredient_id),
  FOREIGN KEY(nabor_id) REFERENCES nabor(id),
  FOREIGN KEY(ingredient_id) REFERENCES ingredient(id)
);

CREATE TABLE IF NOT EXISTS nabor_tech_cart (
  nabor_id bigserial,
  tech_cart_id bigserial,
  PRIMARY KEY(nabor_id, tech_cart_id),
  FOREIGN KEY(nabor_id) REFERENCES nabor(id),
  FOREIGN KEY(tech_cart_id) REFERENCES tech_cart(id)
);

CREATE TABLE IF NOT EXISTS dealer (
  id bigserial,
  name text NOT NULL,
  address text NOT NULL,
  phone text NOT NULL,
  comment text NOT NULL,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS schet (
  id bigserial,
  name text NOT NULL,
  currency text NOT NULL,
  type text NOT NULL,
  start_balance float NOT NULL,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);


CREATE TABLE IF NOT EXISTS postavka (
  id bigserial,
  dealer_id bigserial,
  sklad_id bigserial,
  schet_id bigserial,
  time timestamp NOT NULL,
  deleted boolean DEFAULT FALSE,
  PRIMARY KEY(id)
);

INSERT INTO category_tovar VALUES(DEFAULT, 'Главный экран','image.png');
INSERT INTO category_tovar VALUES(DEFAULT, 'Десерты','image.png');
INSERT INTO category_tovar VALUES(DEFAULT, 'Еда','image.png');
INSERT INTO category_tovar VALUES(DEFAULT, 'Напитки','image.png');
INSERT INTO category_ingredient VALUES(DEFAULT, 'Напитка и вода','image.png');
INSERT INTO category_ingredient VALUES(DEFAULT, 'Сиропы','image.png');
INSERT INTO category_ingredient VALUES(DEFAULT, 'Еда и шоколадки','image.png');
INSERT INTO category_ingredient VALUES(DEFAULT, 'Без категории','image.png');

INSERT INTO tovar VALUES(DEFAULT,'Албени', 1, 'albeni.jpg', 'Фискальный налог', 'шт.', 166.42, 250.0, 83.58, 50);
INSERT INTO tovar VALUES(DEFAULT,'Барни Чоко пай', 1,'barniChoko.jpg', 'Фискальный налог', 'шт.', 107.42, 150.0, 42.58, 40);
INSERT INTO tovar VALUES(DEFAULT,'Боул с курицей', 2,'boulTofu.jpeg', 'Фискальный налог', 'шт.', 262.0, 1890.0, 1628.0, 621);
INSERT INTO tovar VALUES(DEFAULT,'Боул с тофу', 2,'boulTofu.jpeg', 'Фискальный налог', 'шт.', 1250.0, 1890.0, 640, 51);
INSERT INTO tovar VALUES(DEFAULT,'Кола', 3,'image.png', 'cola.jpeg', 'шт.', 262.0, 250.0, -12.0, -5);
INSERT INTO tovar VALUES(DEFAULT,'Сертификат 10.000', 4,'image.png', 'Фискальный налог', 'шт.', 10000.0, 10000.0, 0, 0);

INSERT INTO tech_cart VALUES(DEFAULT,'CHERRY ENERGY', 2,'image.png', 'Фискальный налог', 'шт.', 428.0, 990.0, 562, 131);
INSERT INTO tech_cart VALUES(DEFAULT,'Ice latte с сиропом', 2,'image.png', 'Фискальный налог', 'шт.', 428.0, 990.0, 562, 131);
INSERT INTO tech_cart VALUES(DEFAULT,'Ice latte с эссенцией', 2,'image.png', 'Фискальный налог', 'шт.', 428.0, 990.0, 562, 131);



INSERT INTO sklad VALUES(DEFAULT,'Иманова', 'Иманова 3/1в');
INSERT INTO sklad VALUES(DEFAULT,'Астана Молл', 'Проспект Тәуелсіздік, 34');

INSERT INTO ingredient VALUES(DEFAULT,'19л вода', 1,'image.png', 'шт.', 158.50);
INSERT INTO ingredient VALUES(DEFAULT,'Айс Спаниш эссенция', 2,'image.png', 'кг', 113.29);
INSERT INTO ingredient VALUES(DEFAULT,'Бадьян Корица Гвоздика', 4,'image.png', 'кг', 9992.9);
INSERT INTO ingredient VALUES(DEFAULT,'Какао', 2,'image.png', 'кг', 30.57);
INSERT INTO ingredient VALUES(DEFAULT,'Сок вишневый', 1,'image.png', 'л', 611.35);
INSERT INTO ingredient VALUES(DEFAULT,'Зерна кофе', 1,'image.png', 'кг', 8344.9);
INSERT INTO ingredient VALUES(DEFAULT,'Стаканы летние', 1,'image.png', 'шт', 49.36);
INSERT INTO ingredient VALUES(DEFAULT,'Крышки летние', 1,'image.png', 'шт', 13.0);
INSERT INTO ingredient VALUES(DEFAULT,'Трубочки', 1,'image.png', 'шт', 3.0);
INSERT INTO ingredient VALUES(DEFAULT,'Молоко 2,5%', 1,'image.png', 'л', 520.63);
INSERT INTO ingredient VALUES(DEFAULT,'Сироп на выбор ингридиент', 1,'image.png', 'л', 1 753.17);
INSERT INTO ingredient VALUES(DEFAULT,'Эссенция без сахара', 1,'image.png', 'шт ', 51.86);




INSERT INTO tech_cart_ingredient VALUES(1,5,0.32);
INSERT INTO tech_cart_ingredient VALUES(1,6,0.02);
INSERT INTO tech_cart_ingredient VALUES(1,7,1);
INSERT INTO tech_cart_ingredient VALUES(1,8,1);
INSERT INTO tech_cart_ingredient VALUES(1,9,1);


INSERT INTO tech_cart_ingredient VALUES(2,5,0.32);
INSERT INTO tech_cart_ingredient VALUES(2,6,0.011);
INSERT INTO tech_cart_ingredient VALUES(2,7,1);
INSERT INTO tech_cart_ingredient VALUES(2,8,1);
INSERT INTO tech_cart_ingredient VALUES(2,9,1);
INSERT INTO tech_cart_ingredient VALUES(2,10,0.3);
INSERT INTO tech_cart_ingredient VALUES(2,11,0.03);

INSERT INTO tech_cart_ingredient VALUES(3,5,0.32);
INSERT INTO tech_cart_ingredient VALUES(3,6,0.011);
INSERT INTO tech_cart_ingredient VALUES(3,7,1);
INSERT INTO tech_cart_ingredient VALUES(3,8,1);
INSERT INTO tech_cart_ingredient VALUES(3,9,1);
INSERT INTO tech_cart_ingredient VALUES(3,10,0.3);
INSERT INTO tech_cart_ingredient VALUES(3,11,0.03);
INSERT INTO tech_cart_ingredient VALUES(3,12,2);

INSERT INTO sklad_tovar VALUES(1, 1, 1000);
INSERT INTO sklad_tovar VALUES(1, 2, 580);
INSERT INTO sklad_ingredient VALUES(1, 1, 1000);
INSERT INTO sklad_ingredient VALUES(1, 2, 580);

INSERT INTO shifts VALUES(3, '2022-12-21 02:54:00.606436+00'::date, '2022-12-21 16:17:00.606436+00'::date, 108020, 11870, ,214550,17680.05,198887.88, 3000)
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "openShift", '2022-12-21 02:54:00.606436+00'::date, 108020, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "postavka", '2022-12-21 03:49:00.606436+00'::date, 76450, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "postavka", '2022-12-21 04:24:00.606436+00'::date, 12500, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "postavka", '2022-12-21 04:26:00.606436+00'::date, 10800, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "postavka", '2022-12-21 06:24:00.606436+00'::date, 16800, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "postavka", '2022-12-21 06:38:00.606436+00'::date, 85500, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "postavka", '2022-12-21 08:09:00.606436+00'::date, 12500, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "collection", '2022-12-21 16:15:00.606436+00'::date, 3000, "")
INSERT INTO transactions VALUES(DEFAULT, 3, 6, "closeShift", '2022-12-21 16:17:00.606436+00'::date, 11870, "")



