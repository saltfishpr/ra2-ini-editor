import { Button, ConfigProvider, Layout, message, Typography } from "antd";
import { useEffect, useState } from "react";
import {
  DeleteUnit,
  GetUnit,
  ListAllUnits,
  NextUnitID,
  Open,
  Save,
  SaveUnit,
  UserRules,
} from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";
import "./App.css";
import EmptyState from "./components/EmptyState";
import Toolbar from "./components/Toolbar";
import UnitDetail from "./components/UnitDetail";
import UnitList from "./components/UnitList";

const { Header, Sider, Content } = Layout;

function App() {
  const [units, setUnits] = useState<main.Unit[]>([]);
  const [selectedUnit, setSelectedUnit] = useState<main.Unit | null>(null);

  useEffect(() => {
    ListAllUnits()
      .then((units) => {
        setUnits(units);
      })
      .catch((err) => {
        console.error("Error listing units:", err);
      });
  }, []);

  const handleAddUnit = (type: string, name: string) => {
    NextUnitID(type)
      .then((id) => {
        const newUnit = main.Unit.createFrom({
          type: type,
          id: id,
          name: name,
          properties: [],
        });
        SaveUnit(newUnit)
          .then(() => {
            setSelectedUnit(newUnit);
            ListAllUnits()
              .then((units) => {
                setUnits(units);
              })
              .catch((err) => {
                console.error("Error listing units after add:", err);
              });
          })
          .catch((err) => {
            console.error("Error saving new unit:", err);
          });
      })
      .catch((err) => {
        console.error("Error getting next unit ID:", err);
      });
  };

  const handleSaveUnit = () => {
    if (!selectedUnit) return;

    SaveUnit(selectedUnit)
      .then(() => {
        GetUnit(selectedUnit.type, selectedUnit.id)
          .then((unit) => {
            setSelectedUnit(unit);
          })
          .catch((err) => {
            console.error("Error getting unit after save:", err);
          });
      })
      .catch((err) => {
        console.error("Error saving unit:", err);
      });
  };

  const handleDeleteUnit = () => {
    if (!selectedUnit || !selectedUnit.id || !selectedUnit.type) return;

    DeleteUnit(selectedUnit.type, selectedUnit.id)
      .then(() => {
        ListAllUnits()
          .then((units) => {
            setUnits(units);
            setSelectedUnit(null);
          })
          .catch((err) => {
            console.error("Error listing units after delete:", err);
          });
      })
      .catch((err) => {
        message.error(`删除失败: ${err}`);
      });
  };

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
          <div style={{ display: "flex", gap: 8 }}>
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
            <Button
              onClick={() => {
                Save()
                  .then(() => {
                    console.log("File saved");
                  })
                  .catch((err) => {
                    console.error("Error saving file:", err);
                  });
              }}
            >
              保存文件
            </Button>
            {false && (
              <Button
                onClick={() => {
                  UserRules()
                    .then((rules) => {
                      message.info(`User rules: ${rules}`);
                    })
                    .catch((err) => {
                      console.error("Error getting user rules:", err);
                    });
                }}
              >
                Debug
              </Button>
            )}
          </div>
        </Header>
        <Layout className="content-layout">
          <Sider className="app-sider">
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
            <div style={{ marginBottom: 16 }}>
              <Toolbar
                onAdd={handleAddUnit}
                onSave={handleSaveUnit}
                onDelete={handleDeleteUnit}
                canSave={!!selectedUnit}
                canDelete={!!selectedUnit}
              />
            </div>
            {selectedUnit ? (
              <UnitDetail
                key={`${selectedUnit.type}-${selectedUnit.id}`}
                unit={selectedUnit}
              />
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
