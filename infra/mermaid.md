graph TB
    subgraph User["User"]
        HD[Hardware Device<br>Raspberry Pi]
        Phone[Smartphone]
    end

    subgraph Server["Server Infrastructure"]
        NGINX[NGINX Reverse Proxy]
        NextJS[Next.js Frontend]
        GO[Go Application]
        DB[(PostgreSQL<br>with PostGIS)]
    end

    subgraph External["External"]
        WIFI{WiFi Hotspots}
    end

    %% Connections
    HD --"Connects to"--> WIFI
    HD --"Sends data<br>(TCP/8888)"--> GO
    Phone --"Accesses website<br>(HTTPS)"--> NGINX
    NGINX --"Proxies requests"--> NextJS
    NGINX --"Proxies API requests"--> GO
    NextJS --"Fetches/Sends data"--> GO
    GO --"Stores/Retrieves data"--> DB

    %% Hardware Device Details
    subgraph Hardware["Hardware Device Details"]
        BASH[Bash Script]
        GOCLIENT[Go Client]
    end

    HD --- Hardware
    BASH --"Runs"--> GOCLIENT

    %% Styling
    classDef server fill:#f9f,stroke:#333,stroke-width:2px;
    classDef external fill:#bbf,stroke:#333,stroke-width:2px;
    classDef hardware fill:#bfb,stroke:#333,stroke-width:2px;
    classDef phone fill:#fdb,stroke:#333,stroke-width:2px;
    class NGINX,NextJS,GO,DB server;
    class WIFI external;
    class HD,BASH,GOCLIENT hardware;
    class Phone phone;