import 'dart:convert';

class CheckModel {
  String? tisCheckUrl;
  String? closedAt;
  String? comment;
  double? cost;
  double? discount;
  double? discountPercent;
  String? feedback;
  int? id;
  String? openedAt;
  String? payment;
  String? status;
  double? sum;
  List<TechCartCheck>? techCartCheck;
  List<TovarCheck>? tovarCheck;
  int? userId;
  int? worker;

  CheckModel(
      {tisCheckUrl,
      closedAt,
      comment,
      cost,
      discount,
      discountPercent,
      feedback,
      id,
      openedAt,
      payment,
      status,
      sum,
      techCartCheck,
      tovarCheck,
      userId,
      worker});

  CheckModel.fromJson(Map<String, dynamic> json) {
    tisCheckUrl = json['tisCheckUrl'];
    closedAt = json['closed_at'];
    comment = json['comment'];
    cost = double.parse(json['cost'].toString());
    discount = double.parse(json['discount'].toString());
    discountPercent = double.parse(json['discount_percent'].toString());
    feedback = json['feedback'];
    id = json['id'];
    openedAt = json['opened_at'];
    payment = json['payment'];
    status = json['status'];
    sum = double.parse(json['sum'].toString());
    if (json['techCartCheck'] != null) {
      techCartCheck = <TechCartCheck>[];
      json['techCartCheck'].forEach((v) {
        techCartCheck!.add(TechCartCheck.fromJson(v));
      });
    }
    if (json['tovarCheck'] != null) {
      tovarCheck = <TovarCheck>[];
      json['tovarCheck'].forEach((v) {
        tovarCheck!.add(TovarCheck.fromJson(v));
      });
    }
    userId = json['user_id'];
    worker = json['worker'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['closed_at'] = closedAt;
    data['comment'] = comment;
    data['cost'] = cost;
    data['discount'] = discount;
    data['discount_percent'] = discountPercent;
    data['feedback'] = feedback;
    data['id'] = id;
    data['opened_at'] = openedAt;
    data['payment'] = payment;
    data['status'] = status;
    data['sum'] = sum;
    if (techCartCheck != null) {
      data['techCartCheck'] = techCartCheck!.map((v) => v.toJson()).toList();
    }
    if (tovarCheck != null) {
      data['tovarCheck'] = tovarCheck!.map((v) => v.toJson()).toList();
    }
    data['user_id'] = userId;
    data['worker'] = worker;
    return data;
  }
}

class TechCartCheck {
  int? checkId;
  String? comments;
  double? cost;
  double? discount;
  int? id;
  String? ingredients;
  String? modificators;
  String? name;
  double? price;
  int? quantity;
  int? techCartId;

  TechCartCheck(
      {checkId,
      comments,
      cost,
      discount,
      id,
      ingredients,
      modificators,
      name,
      price,
      quantity,
      techCartId});

  TechCartCheck.fromJson(Map<String, dynamic> json) {
    checkId = json['check_id'];
    comments = json['comments'];
    cost = double.parse(json['cost'].toString());
    discount = double.parse(json['discount'].toString());
    id = json['id'];
    ingredients = json['ingredients'];
    modificators = json['modificators'];
    name = json['name'];
    price = double.parse(json['price'].toString());
    quantity = json['quantity'];
    techCartId = json['tech_cart_id'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['check_id'] = checkId;
    data['comments'] = comments;
    data['cost'] = cost;
    data['discount'] = discount;
    data['id'] = id;
    data['ingredients'] = ingredients;
    data['modificators'] = modificators;
    data['name'] = name;
    data['price'] = price;
    data['quantity'] = quantity;
    data['tech_cart_id'] = techCartId;
    return data;
  }

  List<Modificator> getModificators() {
    if (modificators == null || modificators!.isEmpty || modificators == "[]") {
      return <Modificator>[];
    }

    var modificatorsList = <Modificator>[];
    jsonDecode(modificators!).forEach((v) {
      modificatorsList.add(Modificator.fromJson(v));
    });
    return modificatorsList;
  }
}

class TovarCheck {
  int? checkId;
  String? comments;
  double? cost;
  double? discount;
  int? id;
  String? modifications;
  String? name;
  double? price;
  int? quantity;
  int? tovarId;

  TovarCheck(
      {checkId,
      comments,
      cost,
      discount,
      id,
      modifications,
      name,
      price,
      quantity,
      tovarId});

  TovarCheck.fromJson(Map<String, dynamic> json) {
    checkId = json['check_id'];
    comments = json['comments'];
    cost = double.parse(json['cost'].toString());
    discount = double.parse(json['discount'].toString());
    id = json['id'];
    modifications = json['modifications'];
    name = json['name'];
    price = double.parse(json['price'].toString());
    quantity = json['quantity'];
    tovarId = json['tovar_id'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['check_id'] = checkId;
    data['comments'] = comments;
    data['cost'] = cost;
    data['discount'] = discount;
    data['id'] = id;
    data['modifications'] = modifications;
    data['name'] = name;
    data['price'] = price;
    data['quantity'] = quantity;
    data['tovar_id'] = tovarId;
    return data;
  }

  // List<Modificator> getModificatorNames() {
  //   if (modifications == null ||
  //       modifications!.isEmpty ||
  //       modifications == "[]") {
  //     return <Modificator>[];
  //   }

  //   var modificatorsList = <Modificator>[];
  //   jsonDecode(modifications!).forEach((v) {
  //     modificatorsList.add(Modificator.fromJson(v));
  //   });
  //   return modificatorsList;
  // }
}

class Modificator {
  String? name;
  double? price;

  Modificator({name, price});

  Modificator.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    price = double.parse(json['price'].toString());
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['name'] = name;
    data['price'] = price;
    return data;
  }
}

class CheckItem {
  String name;
  List<String> dopItems;
  String? comment;
  double price;
  int count;

  CheckItem(this.name, this.dopItems, this.comment, this.price, this.count);
}
