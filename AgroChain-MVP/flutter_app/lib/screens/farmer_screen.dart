import 'package:flutter/material.dart';
import 'package:flutter_app/api_service.dart';
import 'package:flutter_app/screens/consumer_screen.dart';

class FarmerScreen extends StatefulWidget {
  const FarmerScreen({super.key});
  @override
  State<FarmerScreen> createState() => _FarmerScreenState();
}

class _FarmerScreenState extends State<FarmerScreen> {
  final _cropController = TextEditingController();
  final _quantityController = TextEditingController();
  final ApiService _apiService = ApiService();
  String _statusMessage = '';

  void _submitProduce() async {
    if (_cropController.text.isEmpty || _quantityController.text.isEmpty) {
      setState(() => _statusMessage = 'Please fill all fields');
      return;
    }
    setState(() => _statusMessage = 'Submitting...');
    try {
      final response = await _apiService.createProduce(
        _cropController.text,
        int.parse(_quantityController.text),
        'Farmer Ramesh', // Hardcoded for MVP
      );
      setState(() {
        _statusMessage = 'Success! Produce ID: ${response['produceId']}';
        _cropController.clear();
        _quantityController.clear();
      });
    } catch (e) {
      setState(() => _statusMessage = 'Error: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('AgriTrace - Farmer Portal'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            tooltip: 'Trace Produce',
            onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (context) => const ConsumerScreen())),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text('Register New Produce', style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 20),
            TextField(controller: _cropController, decoration: const InputDecoration(labelText: 'Crop Name (e.g., Tomatoes)')),
            const SizedBox(height: 10),
            TextField(controller: _quantityController, decoration: const InputDecoration(labelText: 'Quantity (in kg)'), keyboardType: TextInputType.number),
            const SizedBox(height: 20),
            ElevatedButton(onPressed: _submitProduce, child: const Text('Submit to Ledger')),
            const SizedBox(height: 20),
            if (_statusMessage.isNotEmpty) Text(_statusMessage, style: const TextStyle(fontSize: 16), textAlign: TextAlign.center),
          ],
        ),
      ),
    );
  }
}