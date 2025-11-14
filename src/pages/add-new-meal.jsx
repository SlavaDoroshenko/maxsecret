import {
  Panel,
  Grid,
  Container,
  Flex,
  Avatar,
  Typography,
  Button,
  CellList,
  CellAction,
} from "@maxhub/max-ui";
import TabBar from "../components/tabs";
import CircularProgress from "../components/circular-progress";
import { currentNorm, dailyNorm, dailyRemain } from "../mocks/variables";
import ArrowRight from "../assets/arrow-right.svg?react";
import NormStatic from "../components/norm-static";
import MacronutrientBar from "../components/progress-bars";
import Plus from "../assets/plus.svg?react";
import BreakFastIcon from "../assets/breakfast.svg?react";
import MealItem from "../components/meal-item";
import { useMaxBridge } from "../contexts/maxBridgeContext";

const AddNewMeal = ({ meal, date }) => {
  const webApp = useMaxBridge();

  webApp.BackButton.show();
  return (
    <Panel mode="secondary" className="flex flex-col">
      <div className="min-h-dvh bg-inherit pb-[70px]">
        <CellList
          mode="full-width"
          filled={true}
          className=" rounded-b-4xl overflow-hidden "
        >
          <Flex direction="column" align="center" className="p-2">
            <Typography.Headline variant="large-strong" className="mt-2">
              {meal}
            </Typography.Headline>

            <Typography.Title variant="custom" className="text-xl">
              {date}
            </Typography.Title>
          </Flex>
        </CellList>

        <TabBar />
      </div>
    </Panel>
  );
};

export default AddNewMeal;
