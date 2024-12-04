// Device Part Class
export class DevicePart {
    constructor(id, partId, deviceId, name, price, count) {
        this.id = id;
        this.partId = partId;
        this.deviceId = deviceId;
        this.name = name;
        this.price = parseInt(price);
        this.count = parseFloat(count);
    }
}