// @Libs
import React, { useCallback, useMemo, useState } from 'react';
import { message } from 'antd';
import { useHistory } from 'react-router-dom';
import { snakeCase } from 'lodash';
// @Components
import { SourceForm } from './SourceForm';
// @Services
import ApplicationServices from '@service/ApplicationServices';
// @Routes
import { routes } from '@page/SourcesPage/routes';
// @Types
import { FormWrapProps } from '@page/SourcesPage/partials/_common/SourceForm/SourceForm.types';
// @Utils
import { makeObjectFromFieldsValues } from '@util/Form';
import { CollectionSourceData } from '@page/SourcesPage/SourcesPage.types';
import { sourceFormCleanFunctions } from '@page/SourcesPage/partials/_common/SourceForm/sourceFormCleanFunctions';

const SourceFormWrap = ({
  sources = [],
  connectorSource,
  projectId,
  sourceData = {} as SourceData,
  formMode = 'create',
  setSources
}: FormWrapProps) => {
  const history = useHistory();

  const [isPending, switchPending] = useState(false);

  const [isConnected, setConnected] = useState(sourceData.connected);

  const services = useMemo(() => ApplicationServices.get(), []);

  const handleFinish = useCallback(
    async({ collections, ...rest }: SourceData) => {
      switchPending(true);

      const createdSourceData: SourceData = {
        sourceType: sourceFormCleanFunctions.getSourceType(connectorSource),
        sourceProtoType: snakeCase(connectorSource.id),
        ...makeObjectFromFieldsValues<Pick<SourceData, 'config' | 'destinations' | 'sourceId'>>(rest),
        collections: [] as CollectionSource[],
        connected: isConnected
      };

      if (!createdSourceData.connected) {
        const { config, sourceId } = createdSourceData;
        createdSourceData.connected = await sourceFormCleanFunctions.testConnection({ config, sourceId }, connectorSource);
      }

      if (collections) {
        createdSourceData.collections = collections.map((collection: any) => ({
          name: collection.name,
          type: collection.type,
          schedule: collection.schedule,
          parameters: connectorSource.collectionParameters.reduce((accumulator: any, current: any) => {
            return {
              ...accumulator,
              [current.id]: collection[current.id]
            };
          }, {})
        }));
      }

      const payload: CollectionSourceData = {
        sources: formMode === 'edit'
          ? sources.reduce((accumulator: SourceData[], current: SourceData) => [
            ...accumulator,
            current.sourceId !== rest.sourceId
              ? current
              : createdSourceData
          ], [])
          : [...sources, createdSourceData]
      };

      try {
        await services.storageService.save('sources', payload, projectId);

        setSources(payload);

        message.success('New source has been added!');

        history.push(routes.root);

      } catch(error) {
        message.error('Something goes wrong, source hasn\'t been added');
      } finally {
        switchPending(false)
      }
    },
    [isConnected, connectorSource, services.storageService, projectId, sources, history, setSources, formMode]
  );

  return (
    <div className="add-source flex flex-col items-stretch">

      <SourceForm
        formMode={formMode}
        initialValues={sourceData}
        connectorSource={connectorSource}
        isRequestPending={isPending}
        handleFinish={handleFinish}
        sources={sources}
        setConnected={setConnected}
      />
    </div>
  );
};

SourceFormWrap.displayName = 'SourceFormWrap';

export { SourceFormWrap };
