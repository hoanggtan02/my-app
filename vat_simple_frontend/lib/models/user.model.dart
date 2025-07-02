class User {
  final String id;
  final String email;
  final String companyName;

  User({
    required this.id,
    required this.email,
    required this.companyName,
  });

  // Hàm để tạo một đối tượng User từ JSON
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      email: json['email'] ?? '',
      companyName: json['company_name'] ?? '',
    );
  }
}