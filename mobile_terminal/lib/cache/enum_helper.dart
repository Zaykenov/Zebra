class EnumHelper {
  static String enumToStr(dynamic value) {
    return value.toString().substring(value.toString().indexOf('.') + 1);
  }
}
