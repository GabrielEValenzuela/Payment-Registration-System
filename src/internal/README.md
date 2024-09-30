# Internal Folder

> [!NOTE]
> Delete once project has finished

In accordance with the **Screaming Architecture** principle, we organize the project by **business domain** rather than by technical concerns. Each core entity in the payment registration system has its own folder within the `internal/` directory. This structure ensures that the business logic is at the forefront, making it clear what the system **does**, rather than how it is implemented.

For each entity, we will create the following Go files:

### 1. **interface.go**

- **Purpose**: Define the interface(s) that represent the contract for interacting with the entity. This ensures loose coupling between the core business logic and infrastructure.

### 2. **query.go**

- **Purpose**: Implement database queries specific to the entity. These functions will interact with the database and provide the data needed for the business logic.

### 3. **repository.go**

- **Purpose**: Provide the implementation of the repository interface defined in `interface.go`. This is where the actual interaction with the database occurs.

### 4. **service.go**

- **Purpose**: Contains the business logic for the entity. The service uses the repository to interact with the database and applies any business rules before returning data to the caller.

### Example Folder Structure:

```
internal/
│
├── bank/
│   ├── interface.go
│   ├── interface_test.go
│   ├── query.go
│   ├── query_test.go
│   ├── repository.go
│   ├── repository_test.go
│   └── service.go
│   └── service_test.go
│
└── customer/
    ├── interface.go
    ├── interface_test.go
    ├── query.go
    ├── query_test.go
    ├── repository.go
    ├── repository_test.go
    └── service.go
    └── service_test.go

```

### Additional Notes:

- This structure ensures each entity is well-encapsulated and adheres to **single responsibility** principles.
- You can easily extend this architecture for other entities.

By following this structure, we keep the domain logic clean, modular, and easily maintainable, allowing for future changes or enhancements to be implemented with minimal impact on existing code.
