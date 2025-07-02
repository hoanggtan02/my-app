import 'dart:async';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:vat_simple_frontend/api/auth_service.api.dart';
import 'package:vat_simple_frontend/models/user.model.dart'; 

class AuthProvider with ChangeNotifier {
  String? _token;
  User? _user; 
  final AuthService _authService = AuthService();

  bool get isAuth => _token != null;
  String? get token => _token;
  User? get user => _user; 

  Future<void> login(String email, String password) async {
    final token = await _authService.loginUser(email, password);
    _token = token;
    
    // Sau khi đăng nhập, lấy thông tin user
    await fetchAndSetUser();

    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('jwt_token', token);
    notifyListeners();
  }
  
  // Hàm mới để lấy và lưu thông tin user
  Future<void> fetchAndSetUser() async {
    if (_token == null) return;
    try {
      _user = await _authService.getUserProfile(_token!);
      notifyListeners();
    } catch (error) {
      print(error);
      // Có thể xử lý lỗi ở đây, ví dụ logout nếu token không hợp lệ
    }
  }

  Future<bool> tryAutoLogin() async {
    final prefs = await SharedPreferences.getInstance();
    if (!prefs.containsKey('jwt_token')) {
      return false;
    }
    _token = prefs.getString('jwt_token');
    
    // Khi tự động đăng nhập, cũng lấy thông tin user
    await fetchAndSetUser();
    
    notifyListeners();
    return true;
  }

  Future<void> logout() async {
    _token = null;
    _user = null; // <-- Xóa thông tin user khi logout
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('jwt_token');
    notifyListeners();
  }
}