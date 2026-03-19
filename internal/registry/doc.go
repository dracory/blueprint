// Package registry provides the central dependency registry for Dracory.
//
// In Dracory, the Registry is the single source of truth for all core
// dependencies and services used throughout the system. It is responsible
// for:
//
//   1. Providing typed access to core services:
//        - Database connections (GetDatabase())
//        - Logger (GetLogger())
//        - Router (GetRouter())
//        - Config (GetConfig())
//
//   2. Managing static stores that are known at compile time, such as
//      UserStore or OrderStore. These stores may be conditionally initialized
//      based on configuration flags (e.g., enableUserStore).
//
// The Registry is **not** a full application container or harness:
//   - It does **not orchestrate execution** of servers, jobs, or workflows.
//   - It does **not manage dynamic registration** of arbitrary components
//     by string or type (all types are known ahead of time).
//
// Instead, it functions as a **typed dependency registry**:
//   - All services and stores are registered at initialization.
//   - Other components fetch dependencies from the Registry via typed getters.
//   - Conditional initialization allows configuration-driven behavior without
//     exposing unneeded services.
//
// Conceptually, it acts as a central hub for runtime services:
//
//       +-----------------+
//       |     Registry    |
//       |-----------------|
//       | Config          |
//       | DB              |
//       | Logger          |
//       | Router          |
//       | UserStore       |
//       | OrderStore      |
//       +-----------------+
//                ^
//                |
//   Other components fetch dependencies here
//
// This design enforces:
//   - Strong typing at compile time
//   - Centralized dependency management
//   - Configuration-driven conditional initialization
//
// Usage example:
//
//    registry := NewRegistry(config)
//
//    db := registry.GetDatabase()
//    userStore := registry.GetUserStore()  // may be nil if disabled by config
//
package registry
