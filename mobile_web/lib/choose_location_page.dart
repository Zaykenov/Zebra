import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';
import 'package:mobile_web/shoping/src/home/presentation/home_page.dart';
import 'package:permission_handler/permission_handler.dart';

import '../helpers/api_caller.dart';

class ChooseLocation extends StatefulWidget {
  const ChooseLocation({super.key});

  @override
  State<ChooseLocation> createState() => _ChooseLocationState();
}

class _ChooseLocationState extends State<ChooseLocation> {
  LatLng? _currentPosition;
  Set<Marker> zebraLocationMarkers = {};
  @override
  void initState() {
    super.initState();
    _getCurrentLocation();
    startLoadZebraLocations();
  }

  void _getCurrentLocation() async {
    // Check if location permission is granted
    if (await Permission.location.isGranted) {
      Position position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.best,
        timeLimit: const Duration(seconds: 10),
      );

      setState(() {
        _currentPosition = LatLng(position.latitude, position.longitude);
      });
    } else {
      // Request location permission
      var status = await Permission.location.request();
      if (status.isGranted) {
        // If permission is granted, get the current location
        _getCurrentLocation();
      }
    }
  }

  Future<void> startLoadZebraLocations() async {
    try {
      final locations = await ApiCaller().fetchZebraLocations();

      Set<Marker> markers = locations
          .map<Marker>(
            (e) => Marker(
              markerId: MarkerId(e['ShopId'].toString()),
              position: LatLng(e['Latitude'], e['Longitude']),
              infoWindow: InfoWindow(title: e['Name']),
            ),
          )
          .toSet();

      setState(() {
        zebraLocationMarkers = markers;
      });
    } catch (error) {
      // Handle the error as needed
    }
  }

  List<Marker> _sortedZebraLocations() {
    if (_currentPosition == null) return zebraLocationMarkers.toList();

    List<Marker> sortedLocations = zebraLocationMarkers.toList();
    sortedLocations.sort((a, b) {
      double distanceA = _calculateDistance(a.position);
      double distanceB = _calculateDistance(b.position);
      return distanceA.compareTo(distanceB);
    });

    return sortedLocations;
  }

  double _calculateDistance(LatLng location) {
    if (_currentPosition == null) return double.infinity;

    double distance = Geolocator.distanceBetween(
      _currentPosition!.latitude,
      _currentPosition!.longitude,
      location.latitude,
      location.longitude,
    );

    // Convert distance to kilometers
    return distance / 1000;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F7F7),
      appBar: AppBar(title: const Text('Выберите точку')),
      body: _currentPosition == null
          ? const Center(
              child: CircularProgressIndicator(),
            )
          : ListView.builder(
              itemCount: zebraLocationMarkers.length,
              itemBuilder: (context, index) {
                final marker = _sortedZebraLocations()[index];

                return GestureDetector(
                  onTap: () async {
                    await Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (context) => ShopHomePage(
                                  markerId: marker.markerId.value,
                                )));
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Container(
                      decoration: BoxDecoration(
                        border: Border.all(color: Colors.grey), // Add border
                        borderRadius:
                            BorderRadius.circular(8), // Add rounded corners
                      ),
                      child: ListTile(
                        leading: Image.asset('assets/logo.png'),
                        title: Text(marker.infoWindow.title!),
                        subtitle: Text(
                          '${_calculateDistance(marker.position).toStringAsFixed(2)} km from you',
                        ),
                      ),
                    ),
                  ),
                );
              },
            ),
    );
  }
}
