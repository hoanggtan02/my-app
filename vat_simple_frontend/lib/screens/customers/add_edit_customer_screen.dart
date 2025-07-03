import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/api/customer_service.api.dart';
import 'package:vat_simple_frontend/models/customer.model.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';

class AddEditCustomerScreen extends StatefulWidget {
  final Customer? customer;

  const AddEditCustomerScreen({super.key, this.customer});

  @override
  State<AddEditCustomerScreen> createState() => _AddEditCustomerScreenState();
}

class _AddEditCustomerScreenState extends State<AddEditCustomerScreen> {
  final _formKey = GlobalKey<FormState>();
  bool _isLoading = false;

  // Dùng Map để lưu dữ liệu form
  final Map<String, String> _formData = {
    'name': '',
    'tax_code': '',
    'address': '',
    'email': '',
    'phone': '',
  };

  @override
  void initState() {
    super.initState();
    // Nếu là chế độ sửa, điền thông tin cũ vào form
    if (widget.customer != null) {
      _formData['name'] = widget.customer!.name;
      _formData['tax_code'] = widget.customer!.taxCode ?? '';
      _formData['address'] = widget.customer!.address ?? '';
      _formData['email'] = widget.customer!.email ?? '';
      _formData['phone'] = widget.customer!.phone ?? '';
    }
  }

  Future<void> _saveCustomer() async {
    if (!_formKey.currentState!.validate()) return;
    _formKey.currentState!.save();
    
    setState(() => _isLoading = true);
    final token = Provider.of<AuthProvider>(context, listen: false).token!;
    
    try {
      if (widget.customer == null) {
        // --- CHỨC NĂNG TẠO MỚI ---
        await CustomerService().createCustomer(token, _formData);
      } else {
        // --- CHỨC NĂNG CẬP NHẬT ---
        await CustomerService().updateCustomer(token, widget.customer!.id, _formData);
      }

      if (mounted) {
        Navigator.of(context).pop(true); // Trả về true để báo hiệu cần làm mới
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    // Xác định xem đây là màn hình thêm mới hay chỉnh sửa
    final isEditing = widget.customer != null;

    return Scaffold(
      appBar: AppBar(
        title: Text(isEditing ? 'Sửa Khách hàng' : 'Thêm khách hàng mới'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              TextFormField(
                initialValue: _formData['name'],
                decoration: const InputDecoration(labelText: 'Tên khách hàng *'),
                validator: (v) => (v == null || v.isEmpty) ? 'Tên không được để trống' : null,
                onSaved: (v) => _formData['name'] = v!,
              ),
              TextFormField(
                initialValue: _formData['tax_code'],
                decoration: const InputDecoration(labelText: 'Mã số thuế'),
                onSaved: (v) => _formData['tax_code'] = v!,
              ),
              TextFormField(
                initialValue: _formData['address'],
                decoration: const InputDecoration(labelText: 'Địa chỉ'),
                onSaved: (v) => _formData['address'] = v!,
              ),
              TextFormField(
                initialValue: _formData['email'],
                decoration: const InputDecoration(labelText: 'Email'),
                keyboardType: TextInputType.emailAddress,
                onSaved: (v) => _formData['email'] = v!,
              ),
              TextFormField(
                initialValue: _formData['phone'],
                decoration: const InputDecoration(labelText: 'Số điện thoại'),
                keyboardType: TextInputType.phone,
                onSaved: (v) => _formData['phone'] = v!,
              ),
              const SizedBox(height: 20),
              if (_isLoading)
                const Center(child: CircularProgressIndicator())
              else
                ElevatedButton(
                  onPressed: _saveCustomer,
                  child: const Text('Lưu'),
                )
            ],
          ),
        ),
      ),
    );
  }
}