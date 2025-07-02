class Product {
  final String id;
  final String name;
  final String? description;
  final double unitPrice;

  Product({
    required this.id,
    required this.name,
    this.description,
    required this.unitPrice,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      // Chuyển đổi từ số (int hoặc double) sang double
      unitPrice: (json['unit_price'] as num).toDouble(),
    );
  }
}