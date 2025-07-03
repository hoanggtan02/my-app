class Product {
  final String id;
  final String name;
  final String? description;
  final double unitPrice;
  final String? imageUrl;

  Product({
    required this.id,
    required this.name,
    this.description,
    required this.unitPrice,
    this.imageUrl,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    String? finalImageUrl;
    final imageUrlData = json['image_url'];
    if (imageUrlData != null && imageUrlData['Valid'] == true) {
      finalImageUrl = imageUrlData['String'];
    }

    String? finalDescription;
    final descriptionData = json['description'];
    if (descriptionData != null && descriptionData['Valid'] == true) {
      finalDescription = descriptionData['String'];
    }

    return Product(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      description: finalDescription, 
      unitPrice: (json['unit_price'] as num? ?? 0).toDouble(),
      imageUrl: finalImageUrl, 
    );
  }
}