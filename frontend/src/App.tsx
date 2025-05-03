import { Button, ConfigProvider, Layout, Typography, theme } from "antd";
import { useEffect, useState } from "react";
import { GetUnit, ListAllUnits, Open } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";
import "./App.css";
import EmptyState from "./components/EmptyState";
import UnitDetail from "./components/UnitDetail";
import UnitList from "./components/UnitList";

const { Header, Sider, Content } = Layout;

function App() {
  const [units, setUnits] = useState<main.Unit[]>([]);
  const [selectedUnit, setSelectedUnit] = useState<main.Unit | null>(null);
  const { token } = theme.useToken();

  useEffect(() => {
    ListAllUnits()
      .then((units) => {
        setUnits(units);
      })
      .catch((err) => {
        console.error("Error listing units:", err);
      });
  }, []);

  return (
    <ConfigProvider
      theme={{
        components: {
          Collapse: {
            contentPadding: "4px 8px",
          },
        },
      }}
    >
      <Layout className="app-layout">
        <Header className="app-header">
          <Typography.Title level={4} className="app-title">
            RA2 INI 编辑器
          </Typography.Title>
          <Button
            type="primary"
            onClick={() => {
              Open()
                .then(() => {
                  console.log("File opened");
                  ListAllUnits()
                    .then((units) => {
                      setUnits(units);
                    })
                    .catch((err) => {
                      console.error("Error listing units:", err);
                    });
                })
                .catch((err) => {
                  console.error("Error opening file:", err);
                });
            }}
          >
            打开文件
          </Button>
        </Header>
        <Layout className="content-layout">
          <Sider width={280} className="app-sider">
            <UnitList
              units={units}
              onSelectUnit={(unit) => {
                GetUnit(unit.type, unit.id)
                  .then((unit) => {
                    setSelectedUnit(unit);
                  })
                  .catch((err) => {
                    console.error("Error getting unit:", err);
                  });
              }}
              selectedUnit={selectedUnit}
            />
          </Sider>
          <Content className="main-content">
            {selectedUnit ? (
              <UnitDetail unit={selectedUnit} onSave={(unit) => {}} />
            ) : (
              <EmptyState />
            )}
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  );
}

export default App;
