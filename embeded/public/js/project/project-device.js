// Project Device Class
export class ProjectDevice {
    constructor(id, deviceId, projectId, name, price, count, device) {
        this.id = id;
        this.deviceId = deviceId;
        this.projectId = projectId;
        this.name = name;
        this.count = count;
        this.price = price;
        this.device = device;
    }
}