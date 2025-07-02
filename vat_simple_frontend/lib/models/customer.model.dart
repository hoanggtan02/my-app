class Customer {
  final String id;
  final String name;
  final String? taxCode;
  final String? address;
  final String? email;
  final String? phone;

  Customer({
    required this.id,
    required this.name,
    this.taxCode,
    this.address,
    this.email,
    this.phone,
  });

  factory Customer.fromJson(Map<String, dynamic> json) {
    return Customer(
      id: json['id'],
      name: json['name'],
      taxCode: json['tax_code'],
      address: json['address'],
      email: json['email'],
      phone: json['phone'],
    );
  }
}